package main

import (
	"fmt"
	"keys"
	"structures"
)

func main() {
	fmt.Println("Generating RSA key pair and committing to local storage.")
	keys.GenerateKeys()
	structures.MerkleTreeTest()
}
