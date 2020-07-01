package main

import (
	"fmt"
	//"keys"
	//"structures"
    "node"
)

func main() {
    fmt.Println("Launching Node")

    node := node.Make_node()

    node.Blocksize = 2
    node.Killed = false

    node.Run()

}
