package node


import (
	"fmt"
	"pow"
	"structures"
)


//both the arg and the reply
type Block_request struct{
    Index int//requesting the block for this index
}

type Block_request_reply struct{
    Block structures.Block
    Index int//requesting the block for this index
}



func (n *Node) Request_block(index int, peer int) (bool, structures.Block){


    request := Block_request{}
    request.Index = index


    reply := Block_request_reply{}
    block := structures.Block{}

    reply.Block = block


    //we send an empty block for the other node to fill
    n.peers[peer].Call("Send_block", &request, &reply)

    //other node fills out the block in reply, now we can verify it and add to chain

    if pow.Verify_work(reply.Block.Header){
        fmt.Println("Requested block was verified succesfully!")
        return true, reply.Block

    } else {
        fmt.Println("Returned block did not pass verification or was not known")
        return false, reply.Block
    }

}


func (n *Node) Send_block(arg *Block_request, reply *Block_request_reply){

    fmt.Println("Got a block request for index: ", arg.Index)

    if len(n.Chain) >= arg.Index{
        fmt.Println("I have this block, sending....")
        reply.Block = n.Chain[arg.Index]

    } else {
        fmt.Println("I don't have this block. Returning an empty block instead")
        reply.Block = structures.Block{} //should be invalid when verified
    }

}
