package blockchain

import (
	"crypto/sha256"
	"fmt"
)

type MerkeTree struct {
	RootNode *MerkleNode
}

type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		node.Data = hash[:]
	}

	node.Left = left
	node.Right = right

	return &node

}

func NewMerkleTree(data [][]byte) (*MerkeTree, error) {
	var nodes []MerkleNode

	for _, d := range data {
		node := NewMerkleNode(nil, nil, d)
		nodes = append(nodes, *node)
	}

	if len(nodes) == 0 {
		return nil, fmt.Errorf("%s", "No merkle tree node.")
	}

	for len(nodes) > 1 {
		if len(nodes)%2 != 0 {
			dupNode := nodes[len(nodes)-1]
			nodes = append(nodes, dupNode)
		}

		var level []MerkleNode

		for i := 0; i < len(nodes); i += 2 {
			node := NewMerkleNode(&nodes[i], &nodes[i+1], nil)
			level = append(level, *node)
		}

		nodes = level
	}

	tree := &MerkeTree{&nodes[0]}

	return tree, nil
}
