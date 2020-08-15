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

func RequestLastBlock() Block_request_reply {
	reply := Block_request_reply{}
	return reply
}

func (n *Node) Request_block(index int, peer int) (bool, structures.Block) {

	request := Block_request{}
	request.Index = index

	reply := Block_request_reply{}
	block := structures.Block{}

	reply.Block = block

	//we send an empty block for the other node to fill
	/*
		for z, p := range n.PeerPorts {
			fmt.Println(z)
			fmt.Println(p)
		}
	*/
	fmt.Println("node#:", n.Index)
	n.Call(n.PeerPorts[peer], "Server.Send_block", &request, &reply)

	//other node fills out the block in reply, now we can verify it and add to chain
	if pow.Verify_work(reply.Block.Header) {
		fmt.Println("Requested block was verified succesfully!")
		return true, reply.Block

	} else {
		fmt.Println("Returned block did not pass verification or was not known")
		return false, reply.Block
	}

}

//we can have the rpc call this
func (n *Node) Recieve_block(block structures.Block) {

	work_valid := pow.Verify_work(block.Header)

	if !work_valid { //ignore it
		fmt.Println("Block is invalid!")
		return
	}

	//other validity conditions here

	//if it passes them all, then we accept it
	n.Chain = append(n.Chain, block)

}

func (n *Node) Broadcast_complete_block(block structures.Block) (bool, structures.Block) {
	valid, new_block := n.Race_complete_block(block)
	if !valid {
		fmt.Println("ERROR: race to complete block failed")
		return valid, new_block
	}
	for i, _ := range n.PeerPorts {
		if i != n.Index {
			request := Complete_block_request{}
			request.Block = new_block
			reply := Complete_block_reply{}
			ok := n.Call(n.PeerPorts[i], "Server.Download_block", &request, &reply)
			if !ok {
				fmt.Println("Failed to replicate block to peer#", i)
			}
		}
	}
	return valid, new_block
}

/*
   This function simulates broadcasting a block for ever peer to try to mine first.
   The peers race to complete the block. The first peer to complete the block, and gets it verified by the block's owner, wins the race.
   The function will return after the first peer to mine the block, so the other Go threads trying to mine the block will stop.
   The peer who mines the block is recorded, so they can have a reward of some type at a later point.
*/
func (n *Node) Race_complete_block(block structures.Block) (bool, structures.Block) {
	c := make(chan Complete_block_reply)
	for i := 0; i < len(n.PeerPorts); i++ {
		go func(i int) {
			reply := Complete_block_reply{}
			reply.Peer = i
			block := pow.Complete_block(block)
			reply.Block = block
			c <- reply
		}(i)
	}
	for i := 0; i < len(n.PeerPorts); i++ {
		completed := <-c
		if pow.Verify_work(completed.Block.Header) {
			V := &Completed{}
			V.Peer = completed.Peer
			V.BlockIndex = block.Index
			n.Peer_completions = append(n.Peer_completions, V)
			go fmt.Println("Peer ", V.Peer, " completed the block.")
			return true, completed.Block
		}
	}
	return false, block
}
