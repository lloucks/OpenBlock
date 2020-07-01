//This will be the node object
//When it is alive it will continously try to build the chain
//while accepting transactions.

package node

import (
	"fmt"
	"pow"
	"structures"
	"time"

	//"strconv"
	"brpc"
	"encoding/binary"
	"encoding/hex"
	"log"
)

type Node struct {
	Chain     []structures.Block //probably should be the merkle tree
	Privkey   []byte             //I'm not sure on what this type should be
	Pubkey    []byte             //same with this one
	Index     int                //What is this for????
	Cur_block structures.Block   //current block to add transactions to

	Blocksize int //How big our blocks will be (in transaction count, for simplicity)

	Block_time     time.Duration //How long we aim for between blocks (seconds)
	Cur_difficulty int           //how many zeros we want (in bits) at front of hash

	Killed bool //So the node knows to kill itself
}

func Make_node() Node {

	node := Node{}

	node.Block_time = 20 * time.Second
	node.Cur_difficulty = 3
	node.Index = 0

	fmt.Println("Made a client node")
	return node
}

//we can have the rpc call this
func (n *Node) Recieve_block(block structures.Block) {

	fmt.Println("Checking block valididty")
	work_valid := pow.Verify_work(block.Header)

	if !work_valid { //ignore it
		fmt.Println("Block is invalid!")
		return
	}

	//other validity conditions here

	fmt.Println("Block is valid")
	//if it passes them all, then we accept it
	block.Index = block.Index + 1
	n.Chain = append(n.Chain, block)

	fmt.Println("Chain is ", n.Chain)
}

//take the current block and try to solve it (done accepting transactions for now)
func (n *Node) GetLastBlock() structures.Block {
	return n.Chain[len(n.Chain)-1]
}

func (n *Node) Generate_block() {
	//TODO

}

func (n *Node) CreateGenesisBlock() {
	block := structures.Block{}

	block.Index = 0
	block.Header.Prev_block_hash = 0000000000000000000000000000000000000000000000000000000000000000
	block.Header.Difficulty = uint32(n.Cur_difficulty)

	block = pow.Complete_block(block)

	n.Recieve_block(block)
}

//You wait for transactions to fill the block,
//THEN you hash the header and increment nonce until complete.
func (n *Node) MakeBlock() structures.Block {
	block := structures.Block{}
	block.Index = n.Index

	fmt.Println("hashing previous block")

	fmt.Println("Length of chain is ", len(n.Chain))
	hexhash, err := hex.DecodeString(pow.GenerateHash(n.Chain[len(n.Chain)-1].Header))

	block.Header.Prev_block_hash = int(binary.BigEndian.Uint32(hexhash))
	if err != nil {
		log.Fatalf("Critical error converting block hash to int: %v\n", err)

	}
	fmt.Println("hashed")

	//difficulty should be the same as last block unless we adjust it

	block.Header.Difficulty = n.GetLastBlock().Header.Difficulty

	return block
}

//Call every X blocks to ensure time between blocks is consistent(ish)
//20 Seconds between blocks for now
//starting difficulty can be 5, then it can be raised or lowered
func (n *Node) Adjust_difficulty() {

	//blocks per min * 2 gives a difficulty check roughly every 2 mins
	//when the duration between blocks is 20 seconds

	adjust_block_count := int((time.Second*60)/n.Block_time) * 2

	// -1 is so we have one extra time to look at
	blocks := n.Chain[len(n.Chain)-adjust_block_count-1:]

	fmt.Printf("Looking at %v blocks for difficulty calc\n", len(blocks))

	var times []time.Time

	for _, block := range blocks {
		times = append(times, block.Header.Timestamp)
	}

	var differences []time.Duration

	for i, t := range times {
		if i == 0 {
			continue
		}
		differences = append(differences, t.Sub(times[i-1]))
	}

	average := time.Duration(0)
	for _, t := range differences {
		average += t
	}

	average = average / time.Duration(len(differences))

	fmt.Printf("Average block time for last %v blocks is %v\n", adjust_block_count, average)
	fmt.Printf("Goal is %v seconds\n", n.Block_time)

	if average-n.Block_time < -time.Duration(time.Second) {
		n.Cur_difficulty++
	} else if average-n.Block_time > time.Duration(time.Second) {
		n.Cur_difficulty--
	}

	fmt.Printf("Changed difficulty to %v\n\n", n.Cur_difficulty)

}

//RPC that other nodes call to send transactions to this node.

//The header must be changed to proper RPC args/reply strandards when we get there
func (n *Node) recieve_transaction(args *brpc.Args, reply *brpc.Reply) {
	//Pull the transaction out of arguments

	t := args.Transaction

	//Validate transaction

	//if valid, append it.
	n.Cur_block.MTree = n.Cur_block.MTree.AddTransaction(t)

}

//As of right now, we will just have a node building it's own chain

//We will need to build on this when it comes to reciving from others.

func (n *Node) Run() {

	//asuming we are starting a brand new chain everytime for now.
	n.CreateGenesisBlock()
	n.Cur_block = n.MakeBlock()

	go n.local_transaction_loop()

	for !n.Killed {
		//if our block is FULL (to be determined when) then we try to complete it and start
		//a new block
		if n.is_cur_block_full() {
			n.Cur_block = pow.Complete_block(n.Cur_block)
			n.Chain = append(n.Chain, n.Cur_block)
			n.Cur_block = n.MakeBlock()
			fmt.Printf("Added a block to the chain\n")
		}

		time.Sleep(time.Millisecond * 50)
		//else we wait for user input to send transactions
	}

}

func (n *Node) is_cur_block_full() bool {
	if len(n.Cur_block.MTree.Leafs) >= n.Blocksize {
		return true
	} else {
		return false
	}
}

//A goroutine that will wait for user input, make a transaction and add it to the current block
func (n *Node) local_transaction_loop() {

	var input string

	for {
		fmt.Println("Enter author Number: ")
		var authorID int
		_, err := fmt.Scanf("%d", &authorID)
		if err != nil {
			log.Fatal("not valid author ID")
		}
		fmt.Println("Enter text: ")
		fmt.Scanln(&input)

		t := structures.CreateTransaction(input, authorID)

		n.Cur_block.MTree = n.Cur_block.MTree.AddTransaction(t)

		fmt.Printf("Added a transaction to block %v\n", len(n.Chain)+1)
	}

}
