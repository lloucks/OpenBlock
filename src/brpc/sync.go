//RPCs in here will involve nodes connecting to peers
//downloading the chain, and letting each other know the status of the chain
package brpc

import (
  "fmt"
  "node"
  "time"
  "crypto/sha256"
  "net/rpc"
)

type brpc_net struct{
	 Nodes []*node.Node
	 SockNames []string
	 N_index int
}




func (n *brpc_net) Get_next() {
  if n.N_index+1 < len(n.Nodes) {
    n.N_index += 1
  }else {
    n.N_index = 0
  }
}

func (n *brpc_net) Get_prev() {
  if len(n.Nodes) > 0 {
    n.N_index -= 1
  }else {
    n.N_index = len(n.Nodes)-1
  }
}

func (n *brpc_net) Add_new_node() *node.Node {
	return n.Node_startup()
}

func (n *brpc_net) Node_startup() *node.Node {
  fmt.Println("Launching Node")
  i := len(n.Nodes)
  node := node.Make_node(i)
  rpc.Register(node)


  node.Blocksize = 2
  node.Killed = false

  node.Cur_difficulty = 15

  go node.Run()

  //wait for genesis block.
  time.Sleep(time.Second * 1)

  n.Nodes = append(n.Nodes, node)
  n.SockNames = append(n.SockNames, node.SockName)
  for _, node_i := range n.Nodes{
		node_i.PeerSocks = n.SockNames
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

func Make_brpc_network() brpc_net{

    rpc.HandleHTTP()
    network := brpc_net{}

    network.N_index = 0

    return network

}
