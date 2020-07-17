package main

import (
	"fmt"
	"keys"
	"structures"
	"node"
    // "../brpc"
)

func main() {
	fmt.Println("Launching Node")

	node := node.Make_node()

	node.Blocksize = 10
	node.Killed = false

	node.Run()

}

// func main(){
    // brpc.Main()
// }