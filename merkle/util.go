package merkle

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"math/big"
	"os"
	"path/filepath"
	"unsafe"
)

type LeafNodes []*LeafNode

type LeafNode string

func (L *LeafNodes) Len() int {
	return len(*L)
}

func (L *LeafNodes) Less(i, j int) bool {
	bi := new(big.Int)
	_, ok := bi.SetString(string(*(*L)[i]), 16)
	if !ok {
		panic(ERRORLEAFNODEFORMAT)
	}
	bj := new(big.Int)
	_, ok = bj.SetString(string(*(*L)[j]), 16)
	if !ok {
		panic(ERRORLEAFNODEFORMAT)
	}
	res := bi.Cmp(bj)
	if res == -1 {
		return true
	}
	return false
}

func (L *LeafNodes) Swap(i, j int) {
	(*L)[i], (*L)[j] = (*L)[j], (*L)[i]
}

func (leafs LeafNodes) isExist(question LeafNode) (bool, error) {
	que := new(big.Int)
	_, ok := que.SetString(string(question), 16)
	if !ok {
		panic(ERRORLEAFNODEFORMAT)
	}
	length := len(leafs)
	if length==0{
		return false,errors.New(LEAFNODESLENGTH)
	}
	lo, hi := 0, length-1
	for lo <= hi {
		m := (lo + hi) >> 1
		i := new(big.Int)
		i.SetString(string(*leafs[m]), 16)
		if res := i.Cmp(que); res == 0 {
			return true, nil
		} else if res == 1 {
			hi = m - 1
		} else if res == -1 {
			lo = m + 1
		}
	}
	return false, nil
}

func (leafs LeafNodes)IsLeafNode(question LeafNode)(int,bool, error) {
	que := new(big.Int)
	qb:=make([]byte,32,32)
	q,err:=hex.DecodeString(string(question))
	if err!=nil{
		panic(err)
	}
	copy(qb[0:5],q[0:5])
	fmt.Println(qb)
	_= que.SetBytes(qb)
	length := len(leafs)
	if length==0{
		return 0,false,errors.New(LEAFNODESLENGTH)
	}
	lo, hi := 0, length-1
	for lo <= hi {
		m := (lo + hi) >> 1
		i := new(big.Int)
		ib:=make([]byte,32,32)
		b,err:=hex.DecodeString(string(*leafs[m]))
		if err!=nil{
			panic(err)
		}
		copy(ib[0:5], b[0:5])
		fmt.Println(ib)
		i.SetBytes(ib)
		if res := i.Cmp(que); res == 0 {
			return m,true, nil
		} else if res == 1 {
			hi = m - 1
		} else if res == -1 {
			lo = m + 1
		}
	}
	return 0,false, nil
}


func (l *LeafNodes) ValidLeafNode() (*LeafNodes,error) {
	var merkleRoot []string
	err:=filepath.Walk(MERKLETREEPATH,
		func(path string, f os.FileInfo, err error) (error) {
			if f == nil {
				return err
			}
			if f.IsDir() {
				if len(f.Name()) == 64||f.Name()==NEXTMERKLETREE {
					merkleRoot = append(merkleRoot, f.Name())
				}
				return err
			}
			return err
		})
	if err!=nil{
		return nil,err
	}
	for _, root := range merkleRoot {
		nodes := new(LeafNodes)
		path := fmt.Sprintf("%s%s", MERKLETREEPATH, root)
		db, err := leveldb.OpenFile(path, nil)
		if err != nil {
			continue
		}
		i := 0
		key := (*[1]byte)(unsafe.Pointer(&i))
		value, err := db.Get(key[:], nil)
		if err != nil&&root!=NEXTMERKLETREE {
			db.Close()
			return nil,err
		}
		var leafNodes = new(LeafNodes)
		err = json.Unmarshal(value, leafNodes)
		if err != nil&&root!=NEXTMERKLETREE {
			db.Close()
			return nil,err
		}
		if len(*leafNodes)==0{
			db.Close()
			continue
		}
		for _, n := range *l {
			ok, err := leafNodes.isExist(*n)
			if err != nil {
				db.Close()
				return nil,err
			}
			if !ok {
				*nodes = append(*nodes, n)
			}
		}
		l=nodes
		db.Close()
	}
	return l,nil
}