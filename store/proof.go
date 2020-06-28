package store

const (
	MERKLE_LEFT = 0
	MERKLE_RIGHT = 1
)

type StoreProof struct{
	Content []byte
	Nodes []MerklePath
}

type MerklePath struct{
	LeftOrRight byte
	hash []byte
}

