package brpc

import (
    "fmt"
    "net"
    "net/http"
	"net/rpc"
	"node"
	"structures"
)

// RPC Server
type Server struct {
	Node *node.Node
	SockName string
    name string
}

func Call(port string, funcName string) {
    client, err := rpc.DialHTTP("tcp", port)
    if err != nil {
        fmt.Println(err)
    }
    var name string
    err = client.Call(funcName, struct{}{}, &name)
    if err != nil {
        fmt.Println(err)
    }
    fmt.Println("Server name on port", port, "is", name)
}

func Serve(name, port string, n *node.Node) {
    serv := rpc.NewServer()
	s := Server{}
	s.name = name
	s.Node = n
    serv.Register(&s)

    // ===== workaround ==========
    oldMux := http.DefaultServeMux
    mux := http.NewServeMux()
    http.DefaultServeMux = mux
    // ===========================

    serv.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

    // ===== workaround ==========
    http.DefaultServeMux = oldMux
    // ===========================

    l, err := net.Listen("tcp", port)
    if err != nil {
        panic(err)
    }
    go http.Serve(l, mux)
}


func Service(name, port string) {
    serv := rpc.NewServer()
	s := Server{}
	s.name = name
    serv.Register(&s)

    // ===== workaround ==========
    oldMux := http.DefaultServeMux
    mux := http.NewServeMux()
    http.DefaultServeMux = mux
    // ===========================

    serv.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

    // ===== workaround ==========
    http.DefaultServeMux = oldMux
    // ===========================

    l, err := net.Listen("tcp", port)
    if err != nil {
        panic(err)
    }
    go http.Serve(l, mux)
}

type Completed struct {
	Peer       int
	BlockIndex int
}

//both the arg and the reply
type Block_request struct {
	Index int //requesting the block for this index
}
type Block_request_reply struct {
	Block structures.Block
	Index int //requesting the block for this index
}

type Complete_block_request struct {
	Block structures.Block //block to be completed
	Index int              //requesting the block for this index
	Peer  int              //The peer that is trying to verify the block
}
type Complete_block_reply struct {
	Block structures.Block
	Index int //requesting the block for this index
	Peer  int //The peer that is trying to verify the block
}


func (s *Server) Name(arg struct{}, ret *string) error {
    *ret = s.name
    return nil
}

func (s *Server) Send_block(arg *Block_request, reply *Block_request_reply) error {

	fmt.Println("Got a block request for index: ", arg.Index)
	fmt.Println("node-1 #:", s.Node.Index)

	if len(s.Node.Chain) >= arg.Index {
		fmt.Println("I have this block, sending....")
		reply.Block = s.Node.Chain[arg.Index]

	} else {
		fmt.Println("I don't have this block. Returning an empty block instead")
		reply.Block = structures.Block{} //should be invalid when verified
	}

    return nil
}





