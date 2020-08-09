package node

import (
	"fmt"
	"pow"
	"structures"
)

type Verification struct {
	Verifier   int
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
	Block    structures.Block //block to be completed
	Index    int              //requesting the block for this index
	Verifier int              //The peer that is trying to verify the block
}
type Complete_block_reply struct {
	Block    structures.Block
	Index    int //requesting the block for this index
	Verifier int //The peer that is trying to verify the block
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
   This function simulates broadcasting a block for ever peer to try to verify first
*/
func (n *Node) Broadcast_complete_block(block structures.Block) structures.Block {
	c := make(chan Complete_block_reply)
	for i := 0; i < len(n.peers)+1; i++ {
		go func(i int) {
			reply := Complete_block_reply{}
			reply.Verifier = i
			block := pow.Complete_block(block)
			reply.Block = block
			c <- reply
		}(i)
	}
	verified := <-c
	V := &Verification{}
	V.Verifier = verified.Verifier
	V.BlockIndex = block.Index
	n.Peer_verifications = append(n.Peer_verifications, V)
	return verified.Block
}
