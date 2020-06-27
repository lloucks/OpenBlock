package structures

import (
	"crypto/rand"
	"fmt"
)

type MerkleNode struct {
	UID        string
	Index      int
	Parent     *MerkleNode
	LeftChild  *MerkleNode
	RightChild *MerkleNode
	HashedData []byte
	leaf       bool
}

/*
	Creates a new node in the for the merkle tree
*/
func CreateMerkleNode(index int, left *MerkleNode, right *MerkleNode,
	hash []byte, leaf bool) *MerkleNode {
	node := MerkleNode{}
	node.Index = index
	node.HashedData = hash
	node.LeftChild = left
	node.RightChild = right
	node.leaf = leaf
	node.UID = generate_uuid()
	return &node
}

/*
	Generates a Unique Identifier for the Merkle Node
*/
func generate_uuid() (uuid string) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}
