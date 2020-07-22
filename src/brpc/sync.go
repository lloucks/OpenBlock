//RPCs in here will involve nodes connecting to peers
//downloading the chain, and letting each other know the status of the chain
package brpc

import (
  "fmt"
  "node"
  "bufio"
  "os"
  "io"
  "strings"
  "time"
  "labrpc"
)

func node_startup() node.Node {
  fmt.Println("Launching Node")

  node := node.Make_node()

  node.Blocksize = 2
  node.Killed = false

  node.Cur_difficulty = 15

  go node.Run()

  //wait for genesis block.
  time.Sleep(time.Second * 1)

  Add_client_end()

  return node
}

func Add_new_node() {
  node := node_startup()
  nodes = append(nodes, node)
}

func List_nodes() {
  for v :=  range(nodes) {
    fmt.Println(v)
  }
}

func Get_next() {
  if n_index+1 < len(nodes) {
    n_index += 1
  }else {
    n_index = 0
  }
}

func Add_client_end(){
  end := net.MakeEnd(len(RPC_ends))
  RPC_ends = append(RPC_ends, end)
}

var nodes []node.Node
var RPC_ends []*labrpc.ClientEnd
var n node.Node
var n_index int
var net *labrpc.Network

func Cli_prompt() {

  net = labrpc.MakeNetwork()
  n = node_startup()
  nodes = append(nodes, n)
  n_index = 0

	reader := bufio.NewReader(os.Stdin) //create a reader to parse input

	options := map[string]func(){
		"list":   n.Print_chain,
		"verify": n.Verify_chain,
		"post":   n.Create_transaction,
		"make node": Add_new_node,
		"list nodes": List_nodes,
		"next node": Get_next,
	}

	//n.Killed is just there in the case we want to kill it from other functions
	for !n.Killed {
    n = nodes[n_index]
    fmt.Println("Currently on node ", n_index)
		fmt.Println("Enter a node command: (list, verify, post)")
    fmt.Println("or an RPC command: (make node, list nodes, next node)")

		command, err := reader.ReadString('\n')

		if err == io.EOF {
			n.Exit()
			return
		}
		//Clear the newline from input
		command = strings.Replace(command, "\n", "", -1)

		fmt.Println()
		//check if our command is valid
		found := false
		for k, v := range options {
			if command == k {
				found = true
				v()
				break
			}
		}
		if !found {
			fmt.Println("Invalid command, please try again.")
			continue
		}

	}

}
