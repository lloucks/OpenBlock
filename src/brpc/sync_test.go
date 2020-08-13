package brpc

import (
    "log"
    "time"
    "fmt"
	"testing"
)

func TestServe(t *testing.T){
    Service("one", ":2001")
    Service("two", ":2002")
    Call(":2001", "Server.Name")
    Call(":2002", "Server.Name")
}

func TestRPC(t *testing.T){

    n := Make_brpc_network()

    n.Add_new_node()
    n.Add_new_node()

    n.List_nodes()


    node1 := n.Nodes[0]
    node2 := n.Nodes[1]

    //fmt.Println("Giving 1 second for nodes to generate the first blocks")
    time.Sleep(time.Second * 1)

    // fmt.Println("Genesis block on node 1:")
    // fmt.Println(node1.Chain[0].To_string())
    // fmt.Println("Genesis block on node 2:")
    // fmt.Println(node2.Chain[0].To_string())

    // fmt.Println("Node 1 requesting genesis block from node 2")
    valid, blockb := node1.Request_block(0, 1)

    fmt.Println("Recieved block:")
    fmt.Printf(blockb.To_string())

    if !valid{
        log.Fatalf("Block was not valid, FAIL")
    }

    //time.Time stores location data so we should use Unix time to more closely
    //test a distributed system.
    if blockb.Header.Timestamp.Unix() != node2.Chain[0].Header.Timestamp.Unix(){
        log.Fatalf("Recieved block timestamp does not match the timestamp on node 2. FAIL")
    }
}

func TestBroadcastRPC(t *testing.T){

    fmt.Println("-----TestBroadcastRPC-----")
    n := Make_brpc_network()

    n.Add_new_node()
    n.Add_new_node()
    n.Add_new_node()

    n.List_nodes()

    //Add two transactions to complete block.
    n.Nodes[0].Create_transaction_from_input("This is my transaction.")
    n.Nodes[0].Create_transaction_from_input("This is my second transaction.")

    //Wait to complete block.
    //Output from Broadcast_complete_block should output
    //which peer completed the block.
    time.Sleep(time.Second * 1)

    if len(n.Nodes[0].Chain) != 2 {
        log.Fatalf("Block was not created")
    }
}

