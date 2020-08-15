//RPCs in here will involve nodes connecting to peers
//downloading the chain, and letting each other know the status of the chain
package brpc

import (
	"crypto/sha256"
	"fmt"
	"math/rand"
	"node"
	"strconv"
	"time"
)

type brpc_net struct {
	Nodes    []*node.Node
	PortNums []string
	N_index  int
}

func (n *brpc_net) Get_next() {
	if n.N_index+1 < len(n.Nodes) {
		n.N_index += 1
	} else {
		n.N_index = 0
	}
}

func (n *brpc_net) Get_prev() {
	if len(n.Nodes) > 0 {
		n.N_index -= 1
	} else {
		n.N_index = len(n.Nodes) - 1
	}
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func (n *brpc_net) Add_new_node() {
	n.Node_startup()
}

func NodeSock() string {
	d := rand.Intn(1000) + 2000
	s := ":"
	//s := "/var/tmp/blockchain-"
	//s += strconv.Itoa(os.Getuid())
	s += strconv.Itoa(d)
	return s
}

func (n *brpc_net) Node_startup() *node.Node {
	nodeName := RandStringRunes(20)
	fmt.Println("Launching Node")
	i := 0
	if len(n.Nodes) > 0 {
		i = n.Nodes[len(n.Nodes)-1].Index + 1
	}
	node := node.Make_node(i)
	//node.PortNum = NodeSock()
	node.Name = nodeName
	Serve(node.Name, node.PortNum, node)

	node.Blocksize = 2
	node.Killed = false

	node.Cur_difficulty = 15

	go node.Run()

	//wait for genesis block.
	time.Sleep(time.Second * 1)

	n.Nodes = append(n.Nodes, node)
	n.PortNums = append(n.PortNums, node.PortNum)
	for _, node_i := range n.Nodes {
		node_i.PeerPorts = n.PortNums
	}

	return node
}

func (n *brpc_net) List_nodes() {
	var result string
	result += fmt.Sprintf("\nNode List")
	result += fmt.Sprintf("-----------------------------------------------------\n")
	for i, n := range n.Nodes {
		result += fmt.Sprintf("Node #%d\n", i)
		result += fmt.Sprintf("Node's Privkey:%v\n", sha256.Sum256((n.Privkey.D.Bytes())))
	}
	result += fmt.Sprintf("-----------------------------------------------------\n\n")

	//did it this way but wanted to return string. It would mess up our
	//commands if we returned anything so must print out here anyway.
	fmt.Printf(result)
}

func Make_brpc_network() brpc_net {

	network := brpc_net{}

	network.N_index = 0

	return network

}
