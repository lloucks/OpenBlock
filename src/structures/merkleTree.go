package structures

import (
	"bytes"
	"crypto/sha256"
	"log"
)

type MerkleTree struct {
	ID    int
	Root  *MerkleNode
	Leafs []*MerkleNode
}

type MerkleTreeList struct {
	Count  int
	Height int
	Nodes  []*MerkleNode
	UIDs   []string
}

/*
	Creates a merkle tree from the transactions
*/
func CreateMerkleTree(ID int, transactions []Transaction) *MerkleTree {
	var nodes []*MerkleNode
	var leafs []*MerkleNode
	nodeCount := 0
	if len(transactions) == 0 {
		log.Panic("transactions array is of length 0")
		return nil
	}
	nodes, leafs, nodeCount = createLeafNodes(nodes, leafs, transactions)
	//all the levels above the leaf nodes can be created from the leaf nodes
	//amountOfLevels := int(math.Log2(float64(len(transactions))))
	/*each iteration of the for loop creates a row in the merkle tree */
	for len(nodes) > 1 {
		prevNodes := nodes
		nodes, nodeCount = createNodeLevel(prevNodes, nodeCount)
	}
	//nodes[0] is the last node created which is the root node
	tree := &MerkleTree{ID: ID, Root: nodes[0]}
	tree.Leafs = leafs
	return tree
}

/*
	Adds a new transaction to the merkle tree
	Returns the new merkle tree
*/
func (m *MerkleTree) AddTransaction(t *Transaction) *MerkleTree {
	var nodes []*MerkleNode
	var leafs []*MerkleNode
	nodeCount := len(m.Leafs)
	hash := t.ID
	node := CreateMerkleNode(nodeCount, nil, nil, hash, true)
	nodeCount++
	nodes = append(nodes, m.Leafs...)
	nodes = append(nodes, node)
	leafs = append(leafs, m.Leafs...)
	leafs = append(leafs, node)
	for len(nodes) > 1 {
		prevNodes := nodes
		nodes, nodeCount = createNodeLevel(prevNodes, nodeCount)
	}
	//nodes[0] is the last node created which is the root node
	tree := &MerkleTree{ID: m.ID, Root: nodes[0]}
	tree.Leafs = leafs
	return tree
}

func createLeafNodes(nodes []*MerkleNode, leafs []*MerkleNode,
	transactions []Transaction) ([]*MerkleNode, []*MerkleNode, int) {
	//create the leaf nodes in the tree
	nodeCount := 0
	for _, t := range transactions {
		//serialize each transaction so it can be hashed
		hash := t.ID
		node := CreateMerkleNode(nodeCount, nil, nil, hash, true)
		nodeCount++
		nodes = append(nodes, node)
		leafs = append(leafs, node)
	}

	return nodes, leafs, nodeCount
}

func createNodeLevel(prevNodes []*MerkleNode, nodeCount int) ([]*MerkleNode, int) {
	//create one level of the non-leaf nodes in the tree
	var row []*MerkleNode
	if len(prevNodes)%2 != 0 {
		prevNodes = append(prevNodes, prevNodes[len(prevNodes)-1])
	}
	for i := 0; i < len(prevNodes); i += 2 {
		newHash := generateNodeHash(prevNodes[i], prevNodes[i+1])
		node := CreateMerkleNode(nodeCount, prevNodes[i], prevNodes[i+1],
			newHash[:], false)
		prevNodes[i].Parent = node
		prevNodes[i].ChildType = 0
		prevNodes[i+1].Parent = node
		prevNodes[i+1].ChildType = 1
		nodeCount++
		row = append(row, node)
	}
	return row, nodeCount
}

/*
	Generates the hash from the two child nodes in the tree
*/
func generateNodeHash(left *MerkleNode, right *MerkleNode) [32]byte {
	combinedHash := append(left.HashedData, right.HashedData...)
	hash := sha256.Sum256(combinedHash) //hash the combined hashes
	return hash
}

//recurse the tree and verify it is up to date by re-calculating the hashes
//
func (n *MerkleNode) recalculateNodeHashes() ([]byte, error) {
	if n.Leaf {
		return n.HashedData, nil
	}
	rightBytes, err := n.RightChild.recalculateNodeHashes()
	if err != nil {
		return nil, err
	}

	leftBytes, err := n.LeftChild.recalculateNodeHashes()
	if err != nil {
		return nil, err
	}

	combinedHash := append(leftBytes, rightBytes...)
	hash := sha256.Sum256(combinedHash) //hash the combined hashes
	HashedData := hash[:]
	return HashedData, nil
}

//TreeDifference verifies the tree hasn't changed its transactions
//
func (m *MerkleTree) TreeDifference() (bool, error) {
	newMerkleRootHash, err := m.Root.recalculateNodeHashes()
	if err != nil {
		return false, err
	}

	if bytes.Compare(m.Root.HashedData, newMerkleRootHash) == 0 {
		return false, nil
	}
	return true, nil
}

//VerifyTransaction find the target transaction and verief the hashes all the way
//to the root hash of the merkle tree
func (m *MerkleTree) VerifyTransaction(t *Transaction) (bool, error) {
	for _, l := range m.Leafs {
		if bytes.Compare(l.HashedData, t.ID) == 0 {
			currentParent := l.Parent
			current := l
			for currentParent != nil {
				rightBytes := currentParent.RightChild.HashedData
				leftBytes := current.HashedData
				if current.ChildType == 1 {
					rightBytes = current.HashedData
					leftBytes = currentParent.LeftChild.HashedData
				}
				combinedHash := append(leftBytes, rightBytes...)
				hash := sha256.Sum256(combinedHash) //hash the combined hashes
				HashedData := hash[:]
				if bytes.Compare(HashedData, currentParent.HashedData) != 0 {
					return false, nil
				}
				currentParent = currentParent.Parent
				current = current.Parent
			}
			return true, nil
		}
	}

	return false, nil
}

func (m *MerkleTree) CheckHeight(n *MerkleNode, depth int) int {
	if n == nil {
		return depth
	}
	right := m.CheckHeight(n.RightChild, depth) + 1
	left := m.CheckHeight(n.LeftChild, depth) + 1

	if left > right {
		return left
	} else {
		return right
	}
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func (m *MerkleTree) CollectNodes(mtl *MerkleTreeList, n *MerkleNode) {
	if n == nil {
		return
	}
	if !contains(mtl.UIDs, n.UID) {
		mtl.UIDs = append(mtl.UIDs, n.UID)
		mtl.Nodes = append(mtl.Nodes, n)
		m.CollectNodes(mtl, n.RightChild)
		m.CollectNodes(mtl, n.LeftChild)
	}
	return

}

func (m *MerkleTree) CheckMerkleTree(n *MerkleNode) MerkleTreeList {
	var mList []*MerkleNode
	var UIDList []string
	height := m.CheckHeight(n, 0)
	mStruct := MerkleTreeList{
		Count:  0,
		Height: height,
		Nodes:  mList,
		UIDs:   UIDList,
	}
	m.CollectNodes(&mStruct, m.Root)
	return mStruct
}
