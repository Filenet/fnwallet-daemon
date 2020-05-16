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



//添加文件到ipfs中以及生成hash答案库
//func TestVerifyHash(t *testing.T){
//	fmt.Println(time.Now().UnixNano())
//	ha:="Qmd5daWVM2MEbrSiKJhziEnJJcHxzFDsR2nj6vQFBWP8Uf" //上一个区块的加密hash,方便测试使用了ipfs存储的hash,实际为sha256
//	sha256,err:=base58.Decode(ha)
//	fmt.Println(sha256)
//	if err!=nil{
//		panic(err)
//	}
//	hash:=hex.EncodeToString(sha256)
//	//找出ipfs的hash,以及ipfsblock的内容
//	ipfsHash,ipfsBlockRaw,err:=GetIpfsHash(hash,5)
//	if err!=nil{
//		panic(err)
//	}
//	fmt.Println(ipfsHash)
//	fmt.Println(ipfsBlockRaw)
//	//验证ipfs的hash是不是找出来的内容
//	is_IpfsHash:=VerifyHash(string(ipfsHash),ipfsBlockRaw)
//	fmt.Println(is_IpfsHash)
//	fmt.Println(time.Now().UnixNano())
//
//}

func TestOther(t *testing.T){

	//batch:=new(leveldb.Batch)
	//batch.Put([]byte("test"),[]byte("hello"))
	//batch.Put([]byte("ceshi"),[]byte("world"))
	//batch.Put([]byte("wsf"),[]byte("大吊逼"))
	//fmt.Println(batch.Len())
	//batch.Delete([]byte("wsf"))
	//fmt.Println(batch.Len())
	//fmt.Printf("%+v\n",batch)
	//fmt.Println("dump",batch.Dump())
	////db.Write(batch,nil)
	//value,err:=db.Get([]byte("wsf"),nil)
	//if err!=nil{
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(string(value))
	////batch.Put()
	////batch.Load()
	db,err:=leveldb.RecoverFile("/levelbatch",nil)
	if err!=nil{
		fmt.Println("reconverfile")
		panic(err)
	}
	//fmt.Println()
	//db,err:=leveldb.OpenFile("/levelbatch",nil)
	//if err!=nil{
	//	fmt.Println("ee")
	//	panic(err)
	//}
	iterator:=db.NewIterator(nil,nil)
	if iterator.Next(){
		fmt.Println("啥玩意",string(iterator.Key()))
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
		// Remember that the contents of the returned slice should not be modified, and
		// only valid until the next call to Next.
		key := iter.Key()
		value := iter.Value()
		fmt.Println(string(key))
		fmt.Println(string(value))
	}
	snapshot,err:=db.GetSnapshot()
	if err!=nil{
		fmt.Println(err)
	}
	err=db.Put([]byte("z"),[]byte("科技"),nil)
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
	//content,err:=ioutil.ReadFile("D:\\2019-9-2.log")
	//if err!=nil{
	//	panic(err)
	//}
	//fmt.Println("测试",string(content))
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
	fmt.Println("大端",uint64B)
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

