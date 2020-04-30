package filenetipfs

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fnv3/test/merkle"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
	"github.com/syndtr/goleveldb/leveldb"
	"io"
)

var (
	IpfsAddPath      = "http://127.0.0.1:5001/api/v0/add"
	IpfsLsPath       = "http://127.0.0.1:5001/api/v0/ls"
	IpfsBlockRawPath = "http://127.0.0.1:5001/api/v0/block/get"
	ContentType      = make(map[uint8]string)
)

func VerifyHash(ipfsHash string, ipfsBlockRaw []byte) bool {
	h1, err := mh.Sum(ipfsBlockRaw, mh.SHA2_256, -1)
	if err != nil {
		panic(err)
	}
	c1 := cid.NewCidV0(h1)
	if c1.String() == ipfsHash {
		return true
	}
	return false
}

func GetIpfsHash(matchHash string, length int) ([]byte, []byte, error) {
	db, err := leveldb.OpenFile("/leveldb", nil)
	if err != nil {
		return nil, nil, err
	}
	defer db.Close()
	Hb, err := hex.DecodeString(matchHash)
	if err != nil {
		return nil, nil, err
	}
	key := sha256.Sum256(Hb[0:length])
	ipfsHash, err := db.Get(key[:], nil)
	if err != nil {
		return nil, nil, err
	}
	url := IpfsBlockRawPath + "?arg=" + string(ipfsHash)
	ipfsBlockRaw, err := IpfsHttpPostJson(url, nil, false)
	if err != nil {
		return nil, nil, err
	}
	if len(ipfsBlockRaw) > 256*1024+14 {
		return nil, nil, errors.New("block to length")
	}
	return ipfsHash, ipfsBlockRaw, nil
}


func SaveFile(file io.Reader, fileName string) error {
	leafNodes, err := SaveFileToIpfs(fileName, file)
	if err != nil {
		panic(err)
	}
	validLeafNodes, err := leafNodes.ValidLeafNode()
	if err != nil {
		return err
	}
	prepareLeafsNodes, err := merkle.NextMerkleTree()
	if err != nil {
		panic(err)
	}
	*prepareLeafsNodes = append(*prepareLeafsNodes, *validLeafNodes...)
	length := len(*prepareLeafsNodes)
	for length/merkle.LEAFLEN >= 1 {
		leafNodes := (*prepareLeafsNodes)[0:merkle.LEAFLEN]
		if length > merkle.LEAFLEN {
			*prepareLeafsNodes = (*prepareLeafsNodes)[merkle.LEAFLEN:]
			length = len(*prepareLeafsNodes)
		} else {
			prepareLeafsNodes = nil
			length = 0
		}
		tree, _, err := merkle.NewMerkleTree(leafNodes)
		if err != nil {
			panic(err)
		}
		err = merkle.SaveMerkleTree(tree)
		if err != nil {
			return err
		}
	}
	if prepareLeafsNodes != nil {
		if err := merkle.SavePrepareLeafNode(prepareLeafsNodes); err != nil {
			return err
		}
	}
	return nil
}
