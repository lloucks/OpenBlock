package main

import (
	"brpc"
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	Cli_prompt()
}

func Cli_prompt() {

	network := brpc.Make_brpc_network()

	//Make first node
	n := network.Node_startup()
	//do setup on first node eg genesis block

	//add it to the network (officially)
	network.Nodes = append(network.Nodes, n)

	reader := bufio.NewReader(os.Stdin) //create a reader to parse input

	n_index := &network.N_index
	//n.Killed is just there in the case we want to kill it from other functions
	for !n.Killed {
		options := map[string]func(){
			"list":              network.Nodes[*n_index].Print_chain,
			"verify":            network.Nodes[*n_index].Verify_chain,
			"post":              network.Nodes[*n_index].Create_transaction,
			"peer completions": network.Nodes[*n_index].Print_peer_completions,
			"make node":         network.Add_new_node,
			"list nodes":        network.List_nodes,
			"next node":         network.Get_next,
			"previous node":     network.Get_prev,
		}
		fmt.Println("\nCurrent Node")
		fmt.Println("-----------------------------------------------------")
		fmt.Printf("Currently on node #%d\n", *n_index)
		fmt.Printf("Current Node Privkey:%v\n",
			sha256.Sum256((network.Nodes[*n_index].Privkey.D.Bytes())))
		fmt.Println("-----------------------------------------------------")
		fmt.Println("")
		fmt.Println("Enter a node command: (list, verify, post, show miner awards)")
		fmt.Println("or an RPC command: (make node, list nodes, next node, previous node)")

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
