package node

import (
	"fmt"
	"pow"
	"structures"
)

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

func (n *Node) Request_block(index int, peer int) (bool, structures.Block) {

	request := Block_request{}
	request.Index = index

	reply := Block_request_reply{}
	block := structures.Block{}

	reply.Block = block

	//we send an empty block for the other node to fill
	n.peers[peer].Call("Send_block", &request, &reply)

	//other node fills out the block in reply, now we can verify it and add to chain

	if pow.Verify_work(reply.Block.Header) {
		fmt.Println("Requested block was verified succesfully!")
		return true, reply.Block

	} else {
		fmt.Println("Returned block did not pass verification or was not known")
		return false, reply.Block
	}
}

func (n *Node) Foo(){
	for peer, _ := range n.peers {
		args := Block_request{}
		reply := Block_request_reply{}
		result := n.peers[peer].Call("Node.Foo_reply", &args, &reply)
		fmt.Println("Sent RPC to ", peer, " result was ", result)
	}
}

func (n *Node) Foo_reply(arg *Block_request, reply *Block_request_reply){
	fmt.Println("I have recieved the RPC.")
}

func (n *Node) Send_block(arg *Block_request, reply *Block_request_reply) {

	fmt.Println("Got a block request for index: ", arg.Index)

	if len(n.Chain) >= arg.Index {
		fmt.Println("I have this block, sending....")
		reply.Block = n.Chain[arg.Index]

	} else {
		fmt.Println("I don't have this block. Returning an empty block instead")
		reply.Block = structures.Block{} //should be invalid when verified
	}
}

/*
   This function simulates broadcasting a block for ever peer to try to mine first.
   The peers race to complete the block. The first peer to complete the block, and gets it verified by the block's owner, wins the race.
   The function will return after the first peer to mine the block, so the other Go threads trying to mine the block will stop.
   The peer who mines the block is recorded, so they can have a reward of some type at a later point.
*/
func (n *Node) Broadcast_complete_block(block structures.Block) (bool, structures.Block) {
	c := make(chan Complete_block_reply)
	for i := 0; i < len(n.peers)+1; i++ {
		go func(i int) {
			reply := Complete_block_reply{}
			reply.Peer = i
			block := pow.Complete_block(block)
			reply.Block = block
			c <- reply
		}(i)
	}
	for i := 0; i < len(n.peers)+1; i++ {
		completed := <-c
		if pow.Verify_work(completed.Block.Header) {
			V := &Completed{}
			V.Peer = completed.Peer
			V.BlockIndex = block.Index
			n.Peer_completions = append(n.Peer_completions, V)
			return true, completed.Block
		}
	}
	return false, block
}
