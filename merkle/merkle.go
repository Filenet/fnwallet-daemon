package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/polydawn/refmt/json"
	"github.com/syndtr/goleveldb/leveldb"
	"go4.org/sort"
	"os"
	"path/filepath"
	"unsafe"
)

const (
	LEAFLEN        = 2
	MERKLETREEPATH = "/leveldb/root/"
	NEXTMERKLETREE = "402051f4be0cc3aad33bcf3ac3d6532b"
	PRV            = 40
)

type MerkleTree []*LeafNodes

func NewMerkleTree(nodes LeafNodes) (*MerkleTree, LeafNode, error) {
	if len(nodes) != LEAFLEN {
		return nil, "", errors.New(ERRORNODENUMBER)
	}
	sort.Sort(&nodes)
	merkleTree := new(MerkleTree)
	NewNodes := nodes
	*merkleTree = append(*merkleTree, &NewNodes)
	for {
		if len(nodes) == 1 {
			break
		}
		NewNodes, err := NewMerkleNode(nodes)
		if err != nil {
			return nil, "", err
		}
		*merkleTree = append(*merkleTree, &NewNodes)
		nodes = NewNodes
	}
	depth := len(*merkleTree)
	root := (*(*merkleTree)[depth-1])[0]
	return merkleTree, *root, nil
}

func NewMerkleNode(nodes LeafNodes) (LeafNodes, error) {
	var newLevel = new(LeafNodes)
	g := len(nodes)
	if g%2 != 0 {
		nodes = append(nodes, nodes[g-1])
	}
	L := len(nodes) / 2
	for i := 0; i < L; i++ {
		left, err := hex.DecodeString(string(*nodes[2*i]))
		if err != nil {
			return nil, err
		}
		right, err := hex.DecodeString(string(*nodes[2*i+1]))
		if err != nil {
			return nil, err
		}
		data := append(left[:], right[:]...)
		newNode := sha256.Sum256(data)
		leafNode := LeafNode(hex.EncodeToString(newNode[:]))
		*newLevel = append(*newLevel, &leafNode)
	}
	return *newLevel, nil
}

func SaveMerkleTree(merkleTree *MerkleTree) error {
	length := len(*merkleTree)
	dbName := (*(*merkleTree)[length-1])[0]
	db, err := leveldb.OpenFile(MERKLETREEPATH+string(*dbName), nil)
	if err != nil {
		return err
	}
	defer db.Close()
	for n, tier := range *merkleTree {
		value, err := json.Marshal(tier)
		if err != nil {
			return err
		}
		key := (*[1]byte)(unsafe.Pointer(&n))
		if err := db.Put(key[:], value, nil); err != nil {
			continue
		}
	}
	return nil
}

func NextMerkleTree() (*LeafNodes, error) {
	db, err := leveldb.OpenFile(MERKLETREEPATH+NEXTMERKLETREE, nil)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var i = 0
	key := (*[1]byte)(unsafe.Pointer(&i))
	value, err := db.Get(key[:], nil)
	if err != nil {
		value = []byte{}
	}
	var leafNodes = new(LeafNodes)
	if len(value) != 0 {
		err = json.Unmarshal(value, leafNodes)
		if err != nil {
			panic(err)
		}
	}
	return leafNodes, nil
}

func SavePrepareLeafNode(l *LeafNodes) error {
	db, err := leveldb.OpenFile(MERKLETREEPATH+NEXTMERKLETREE, nil)
	if err != nil {
		return err
	}
	defer db.Close()

	value, err := json.Marshal(l)
	if err != nil {
		return err
	}
	i := 0
	key := (*[1]byte)(unsafe.Pointer(&i))
	if err := db.Put(key[:], value, nil); err != nil {
		return err
	}

	return nil
}


func SearchLeafNode(random string) (int, bool, error) {
	var merkleRoot []string
	err := filepath.Walk(MERKLETREEPATH,
		func(path string, f os.FileInfo, err error) (error) {
			if f == nil {
				return err
			}
			if f.IsDir() {
				if len(f.Name()) == 64 {
					merkleRoot = append(merkleRoot, f.Name())
				}
				return err
			}
			return err
		})
	if err != nil {
		return 0, false, err
	}
	for _, root := range merkleRoot {
		path := fmt.Sprintf("%s%s", MERKLETREEPATH, root)
		db, err := leveldb.OpenFile(path, nil)
		if err != nil {
			continue
		}
		i := 0
		key := (*[1]byte)(unsafe.Pointer(&i))
		value, err := db.Get(key[:], nil)
		if err != nil {
			db.Close()
			return 0, false, err
		}
		var leafNodes = new(LeafNodes)
		err = json.Unmarshal(value, leafNodes)
		if err != nil {
			db.Close()
			return 0, false, err
		}

		sub, ok, err := leafNodes.IsLeafNode(LeafNode(random))
		if err != nil {
			db.Close()
			return 0, false, err
		}
		if ok {
			return sub, true, nil
		}
		db.Close()
	}
	return 0, false, nil
}
