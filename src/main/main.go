package main

import (
	"fmt"
	"keys"
	"structures"
	"node"
  "brpc"
  "time"
)

func main() {
	fmt.Println("Launching Node")

	node := node.Make_node()

	node.Blocksize = 2
	node.Killed = false

  node.Cur_difficulty = 15
  go node.Run()
  //give it 1 second to start up
  time.Sleep(time.Second * 1)
  fmt.Println("----------------------------------------------------------------------")

	node.Cli_prompt()

}

