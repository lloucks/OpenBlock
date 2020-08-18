package brpc

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"node"
	"strconv"
	"structures"
	"pow"
	"errors"
)

// RPC Server
type Server struct {
	Node    *node.Node
	PortNum string
	name    string
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

	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	fmt.Println("Using port:", l.Addr().(*net.TCPAddr).Port)
	p := ":"
	p += strconv.Itoa(l.Addr().(*net.TCPAddr).Port)
	s.Node.PortNum = p
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
	fmt.Println("node#:", s.Node.Index)

	if len(s.Node.Chain) >= arg.Index {
		fmt.Println("I have this block, sending....")
		reply.Block = s.Node.Chain[arg.Index]

	} else {
		fmt.Println("I don't have this block. Returning an empty block instead")
		reply.Block = structures.Block{} //should be invalid when verified
	}

	return nil
}

func (s *Server) Download_block(arg *Complete_block_request, reply *Complete_block_reply) error {
	fmt.Println("Downloading block with index: ", arg.Block.Index)
	fmt.Println("node#:", s.Node.Index)
	work_valid := pow.Verify_work(arg.Block.Header)
	if !work_valid { //ignore it
		fmt.Println("Block is invalid! Cannot download block for peer:", arg.Peer)
		return errors.New("Block is invalid! Cannot download block for peer")
	}
	s.Node.Chain = append(s.Node.Chain, arg.Block)
	fmt.Println("Downloaded block to chain for node: ", s.Node.Index)
	return nil
}
