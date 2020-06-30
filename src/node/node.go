//This will be the node object
//When it is alive it will continously try to build the chain
//while accepting transactions.


package node

import (
    "structures"
    "pow"
    "time"
    "fmt"
    "strconv"
)




type Node struct{
    Chain []structures.Block //probably should be the merkle tree
    Privkey []byte //I'm not sure on what this type should be
    Pubkey []byte //same with this one
    Index int //What is this for????
    Cur_block structures.Block //current block to add transactions to

    Block_time time.Duration //How long we aim for between blocks (seconds)
    Cur_difficulty int //how many zeros we want (in bits) at front of hash
}



func Make_node() Node{

    node := Node{}

    node.Block_time = 20*time.Second
    node.Cur_difficulty = 3
    node.Index = 0
    return node
}

//we can have the rpc call this
func (n *Node) Recieve_block(block structures.Block){

    work_valid := pow.Verify_work(block.Header)

    if !work_valid{//ignore it
        return
    }

    //other validity conditions here

    //if it passes them all, then we accept it
    block.Index = block.Index+1
    n.Chain = append(n.Chain, block)
}

//take the current block and try to solve it (done accepting transactions for now)
func (n *Node) GetLastBlock() structures.Block{
    return n.Chain[len(n.Chain) -1]
}

func (n *Node) Generate_block(){
//TODO

}



func (n *Node) CreateGenesisBlock(){
    block := structures.Block{}

    block.Index = 0
    block.Header.Prev_block_hash = 0000000000000000000000000000000000000000000000000000000000000000
    block.Header.Difficulty = uint32(n.Cur_difficulty)

    block = pow.Complete_block(block)

    n.Recieve_block(block)
}


//this function makes no sense. You wait for transactions to fill the block,
//THEN you hash the header and increment nonce until complete.
func (n *Node) MakeBlock() structures.Block{
    block := structures.Block{}
    block.Index = n.Index
    var err error
    block.Header.Prev_block_hash, err = strconv.Atoi(pow.GenerateHash(n.Chain[n.Index-1].Header))
    if err != nil{
        fmt.Println("Critical error converting block hash to int")
    }

    //difficulty should be the same as last block unless we adjust it

    block.Header.Difficulty = n.GetLastBlock().Header.Difficulty

    return block
}

//Call every X blocks to ensure time between blocks is consistent(ish)
//20 Seconds between blocks for now
//starting difficulty can be 5, then it can be raised or lowered
func (n *Node) Adjust_difficulty(){

    //blocks per min * 2 gives a difficulty check roughly every 2 mins
    //when the duration between blocks is 20 seconds

    adjust_block_count := int((time.Second*60)/n.Block_time) * 2

    // -1 is so we have one extra time to look at
    blocks := n.Chain[len(n.Chain) - adjust_block_count - 1:]

    fmt.Printf("Looking at %v blocks for difficulty calc\n", len(blocks))

    var times []time.Time

    for _, block := range(blocks){
        times = append(times, block.Header.Timestamp)
    }

    var differences []time.Duration

    for i, t := range(times){
        if i == 0{
            continue
        }
        differences = append(differences, t.Sub(times[i-1]))
    }

    average := time.Duration(0)
    for _, t := range(differences){
        average += t
    }

    average = average/time.Duration(len(differences))

    fmt.Printf("Average block time for last %v blocks is %v\n", adjust_block_count, average)
    fmt.Printf("Goal is %v seconds\n", n.Block_time)

    if average - n.Block_time < -time.Duration(time.Second){
        n.Cur_difficulty ++
    } else if average - n.Block_time > time.Duration(time.Second){
        n.Cur_difficulty --
    }

    fmt.Printf("Changed difficulty to %v\n\n", n.Cur_difficulty)

}
