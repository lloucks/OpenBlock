//RPCs in here will involve nodes connecting to peers
//downloading the chain, and letting each other know the status of the chain
package brpc

import (
  "fmt"
  "node"
  "time"
  "labrpc"
  "crypto/sha256"
)

type brpc_net struct{
     Nodes []*node.Node
     RPC_ends []*labrpc.ClientEnd
     Net *labrpc.Network
     N_index int
}

func (n *brpc_net) Add_client_end(){
  end := n.Net.MakeEnd(len(n.RPC_ends))
  n.RPC_ends = append(n.RPC_ends, end)
}

func (n *brpc_net) Add_server(){
  last_added := len(n.Nodes)-1
  svc := labrpc.MakeService(n.Nodes[last_added])
  svr := labrpc.MakeServer()
  svr.AddService(svc)
  n.Net.AddServer(last_added, svr)

  for index, _ := range n.RPC_ends {
    //n.Net.Connect(n.RPC_ends[index], svr)
    n.Net.Enable(index, true)
  }
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

  n.Add_client_end()

  return node
}

func (n *brpc_net) Add_new_node() {
    node := n.Node_startup()

    for idx, x := range(n.RPC_ends[:len(n.RPC_ends)-1]){
        node.Add_peer(*x)
        n.Nodes[idx].Add_peer(*n.RPC_ends[len(n.RPC_ends)-1])
    }

    n.Nodes = append(n.Nodes, node)

    n.Add_server()
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

    network := brpc_net{}

    network.Net = labrpc.MakeNetwork()

    return network

}
