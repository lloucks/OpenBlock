//RPCs in here will involve nodes connecting to peers
//downloading the chain, and letting each other know the status of the chain
package brpc

import (
  "fmt"
  "node"
  "time"
  "labrpc"
  "crypto/sha256"
  crand "crypto/rand"
  "encoding/base64"
)

type brpc_net struct{
     Nodes []*node.Node
     RPC_ends []*labrpc.ClientEnd
     Net *labrpc.Network
     N_index int
     endnames [][]string
}

//from raft/config.go to generate random string
func randstring(n int) string {
  b := make([]byte, 2*n)
  crand.Read(b)
  s := base64.URLEncoding.EncodeToString(b)
  return s[0:n]
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

func (n *brpc_net) Node_startup() *node.Node {
  fmt.Println("Launching Node")

  node := node.Make_node()

  node.Blocksize = 2
  node.Killed = false

  node.Cur_difficulty = 15

  go node.Run()

  //wait for genesis block.
  time.Sleep(time.Second * 1)

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

func Make(n *node.Node, ends []*labrpc.ClientEnd) *node.Node {
  n.Add_peers(ends)
  return n
}

func (n *brpc_net) Add_new_node() {
    node := n.Node_startup()
    n.Nodes = append(n.Nodes, node)

    //get the index of the newest node.
    node_index := len(n.Nodes)-1

    endnames := make([]string, node_index+1)
    n.endnames = append(n.endnames, endnames)
    for index := range n.endnames[node_index]{
      n.endnames[node_index][index] = randstring(20)
    }

    ends := make([]*labrpc.ClientEnd, node_index+1)
    for index := range ends {
      ends[index] = n.Net.MakeEnd(n.endnames[node_index][index])
      n.Net.Connect(n.endnames[node_index][index], index)
    }

    new_node := Make(node, ends)
    svc := labrpc.MakeService(new_node)
    srv := labrpc.MakeServer()
    srv.AddService(svc)
    n.Net.AddServer(node_index, srv)

    n.Grow_endpoints()
}

func Make_brpc_network() brpc_net{

    network := brpc_net{}

    network.Net = labrpc.MakeNetwork()

    //create the first endname
    network.endnames = make([][]string, 1)

    return network

}
