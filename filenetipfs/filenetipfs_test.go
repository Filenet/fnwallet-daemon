package filenetipfs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
	"github.com/syndtr/goleveldb/leveldb"
	"math/rand"
	"os"
	"testing"
)

func TestOther(t *testing.T){
	db,err:=leveldb.RecoverFile("/levelbatch",nil)
	if err!=nil{
		fmt.Println("reconverfile")
		panic(err)
	}
	iterator:=db.NewIterator(nil,nil)
	if iterator.Next(){
		fmt.Println("db.iterator ",string(iterator.Key()))
	}
	leveldb.RecoverFile("/levelbatch",nil)

}

func TestLvelDB(t *testing.T){
	db,err:=leveldb.OpenFile("/levelbatch",nil)
	if err!=nil{
		panic(err)
	}
	iter:=db.NewIterator(nil,nil)
	fmt.Println(iter)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Println(string(key))
		fmt.Println(string(value))
	}
	snapshot,err:=db.GetSnapshot()
	if err!=nil{
		fmt.Println(err)
	}
	err=db.Put([]byte("z"),[]byte("technology"),nil)
	if err!=nil{
		panic(err)
	}
	value,_:=snapshot.Get([]byte("z"),nil)
	defer snapshot.Release()
	fmt.Println("snapshot",snapshot.String(),value)
}

func TestSeed(t *testing.T){
	rand.Seed(10)
	fmt.Println(rand.Uint64())
}

func TestIpfsCIDV0(t *testing.T){
	file,err:=os.Open("D:\\2019-9-2.log")
	if err!=nil{
		panic(err)
	}
	defer file.Close()
	b:=make([]byte,256*1024,256*1024)
	n,err:=file.Read(b)
	if err!=nil{
		panic(err)
	}
	fmt.Println(n)
	uint64B:=make([]byte,8,8)
	binary.BigEndian.PutUint16(uint64B,90)
	fmt.Println("big-endian",uint64B)
	uint64B1:=make([]byte,8,8)
	binary.BigEndian.PutUint16(uint64B1,0)
	uint64B2:=make([]byte,8,8)
	binary.BigEndian.PutUint16(uint64B2,80)
	newB:=bytes.Join([][]byte{uint64B,uint64B1,uint64B2,b},[]byte{})
	fmt.Println(newB)
	multihash,err:=mh.Sum(newB,mh.SHA2_256,-1)
	if err!=nil{
		panic(err)
	}
	fmt.Println(cid.NewCidV0(multihash))
}

