package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

const (
	LEAFLEN = 1048576
)

type Nodes []*Node

type Node [32]byte

type MerkleTree []*Nodes

func NewMerkleTree(nodes Nodes) (*MerkleTree, string, error) {
	if len(nodes) != LEAFLEN {
		return nil, "", errors.New(ERRORNODENUMBER)
	}
	merkleTree := new(MerkleTree)
	NewNodes:=nodes
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
	root := hex.EncodeToString((*(*merkleTree)[depth-1])[0][:])
	return merkleTree, root, nil
}

func NewMerkleNode(nodes Nodes) (Nodes, error) {
	var newLevel = new(Nodes)
	g := len(nodes)
	if g%2 != 0 {
		nodes = append(nodes, nodes[g-1])
	}
	L := len(nodes) / 2
	for i := 0; i < L; i++ {
		left := nodes[2*i]
		right := nodes[2*i+1]
		data := append(left[:], right[:]...)
		newNode := Node(sha256.Sum256(data))
		*newLevel = append(*newLevel, &newNode)
	}
	return *newLevel, nil
}
