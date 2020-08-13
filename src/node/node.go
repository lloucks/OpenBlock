//This will be the node object
//When it is alive it will continously try to build the chain
//while accepting transactions.

package node

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"pow"
	"strings"
	"structures"
	"time"

	//"strconv"
	"crypto/rsa"
	"keys"
	"log"

  "net/rpc"

)


type Node struct {
	//MTree     *structures.MerkleTree
	Chain     []structures.Block //probably should be the merkle tree
	Privkey   rsa.PrivateKey
	Pubkey    rsa.PublicKey
	Cur_block structures.Block //current block to add transactions to

	Blocksize int //How big our blocks will be (in transaction count, for simplicity)

	Block_time     time.Duration //How long we aim for between blocks (seconds)
	Cur_difficulty int           //how many zeros we want (in bits) at front of hash

	Killed bool //So the node knows to kill itself

	Peer_completions []*Completed
	Index int
	Name string
	PortNum string
	PeerPorts []string
	
	
}

//Structs for RPC calls. Right now they only have block.
type RequestBlockArgs struct{
    Block structures.Block
}

type RequestBlockReply struct{
    Block structures.Block
}








//We will need this function at some point.
//If we want to filter our results for the RPC calls.
func (n *Node) Call(PortNum string, rpcname string, args interface{}, reply interface{}) bool {
	// c, err := rpc.DialHTTP("tcp", "127.0.0.1"+":1234")
	c, err := rpc.DialHTTP("tcp", PortNum)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	defer c.Close()
    //fmt.Println(reply)
	err = c.Call(rpcname, args, reply)
	if err == nil {
		return true
	}

	fmt.Println(err)
	return false
}

func Make_node(i int) *Node {

	node := Node{}
	node.Index = i

	node.Block_time = 20 * time.Second
	node.Cur_difficulty = 3
	tmp_privKey, tmp_pubKey := keys.GetKeys()
	node.Privkey = *tmp_privKey
	node.Pubkey = *tmp_pubKey

	fmt.Println("Made a client node")
    //node.Server() //<-------------------This line makes the node live, and Serve as Server. Ther Server function is defined above. I 
                  // Haven't tested it, but we might need to return a pointer. I may be wrong.  
	return &node
}



//take the current block and try to solve it (done accepting transactions for now)
func (n *Node) GetLastBlock() structures.Block {
	return n.Chain[len(n.Chain)-1]
}

func (n *Node) CreateGenesisBlock() {
	block := structures.Block{}

	block.Index = 0
	block.Header.Prev_block_hash = [32]byte{} //all zeroes by default
	block.Header.Difficulty = uint32(n.Cur_difficulty)

	block = pow.Complete_block(block)

	n.Recieve_block(block)

	fmt.Println("Created genesis block")
}

//You wait for transactions to fill the block,
//THEN you hash the header and increment nonce until complete.
func (n *Node) MakeBlock() structures.Block {
	block := structures.Block{}
	block.Index = len(n.Chain)

	//fmt.Println("hashing previous block")

	//fmt.Println("Length of chain is ", len(n.Chain))
	hash := pow.GenerateHash(n.Chain[len(n.Chain)-1].Header)

	block.Header.Prev_block_hash = hash
	//fmt.Println("hashed previous block header and stored in current block")

	//difficulty should be the same as last block unless we adjust it

	block.Header.Difficulty = n.GetLastBlock().Header.Difficulty

	block.Block_size = uint32(n.Blocksize)

	return block
}

//Call every X blocks to ensure time between blocks is consistent(ish)
//20 Seconds between blocks for now
//starting difficulty can be 5, then it can be raised or lowered
func (n *Node) Adjust_difficulty() {

	//blocks per min * 2 gives a difficulty check roughly every 2 mins
	//when the duration between blocks is 20 seconds

	//flat 5 seems to work better than doing it dynamically
	adjust_block_count := 5 //int((time.Second*60)/n.Block_time) * 2

	//make sure chain is longer than we need to look at otherwise don't adjust
	if len(n.Chain) < adjust_block_count+1 {
		return
	}

	// -1 is so we have one extra time to look at
	blocks := n.Chain[len(n.Chain)-adjust_block_count-1:]

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

	fmt.Printf("Set difficulty to %v\n\n", n.Cur_difficulty)

}

//RPC that other nodes call to send transactions to this node.
//As of right now, we will just have a node building it's own chain
//We will need to build on this when it comes to reciving from others.
//This is a concurrent go_routine
func (n *Node) Run() {

	//asuming we are starting a brand new chain everytime for now.
	n.CreateGenesisBlock()

	n.Cur_block = n.MakeBlock()

	fmt.Println("Node started")

	for !n.Killed {
		//if our block is FULL(transaction count) then we try to complete it and start
		//a new block. We wait until full as there is no monetary incentive for nodes to work on a block.
		//All nodes on the chain are 'lazy', they only work on blocks when they need to.
		if n.is_cur_block_full() {
			n.Cur_block = pow.Complete_block(n.Cur_block)
			n.Chain = append(n.Chain, n.Cur_block)
			fmt.Println("Completed block ", n.Cur_block.Index+1)
			n.Cur_block = n.MakeBlock()
		}

		time.Sleep(time.Millisecond * 50)
		//else we wait for user input to send transactions
	}

}

func (n *Node) is_cur_block_full() bool {
	num_transactions := 0
	if (n.Cur_block.MTree) != nil {
		num_transactions = len(n.Cur_block.MTree.Leafs)
	}
	if num_transactions >= n.Blocksize {
		return true
	} else {
		return false
	}
}


//A goroutine that will wait for user input, make a transaction and add it to the current block
func (n *Node) Create_transaction() {
	var input string
	//fmt.Println("Enter author Number: ")
	var authorID int
	//_, err := fmt.Scanf("%d", &authorID)
	//if err != nil {
	//	log.Fatal("not valid author ID")
	//}

	authorID = 0 //temporary until we figure out how to remove this if not needed

	fmt.Println("Enter text: ")
	reader := bufio.NewReader(os.Stdin)
	//fmt.Scanln(&input)
	input, _ = reader.ReadString('\n')
	// convert CRLF to LF
	input = strings.Replace(input, "\n", "", -1)

	t := structures.CreateTransaction(input, authorID)
	t.Signature = structures.SignTransaction(t)
	if n.Cur_block.MTree == nil {
		var transactions []structures.Transaction
		transactions = append(transactions, *t)
		n.Cur_block.MTree = structures.CreateMerkleTree(1, transactions)
	} else {
		n.Cur_block.MTree = n.Cur_block.MTree.AddTransaction(t)
	}
	n.Chain[len(n.Chain)-1] = n.Cur_block
	//done <- true

	fmt.Printf("Added a transaction to block %v\n", len(n.Chain)+1)
	//fmt.Printf("Amount of leafs in Merkle Tree %v\n", len(n.Cur_block.MTree.Leafs))
}

//Clean up goes here
func (n *Node) Exit() {
	fmt.Println("Quiting.....")
	n.Killed = true

}

func (n *Node) Verify_chain() {
	fmt.Println("Verifying chain...")

}

func (n *Node) Print_chain() {
	fmt.Printf("Number of Blocks %d\n", len(n.Chain))
	totalTrans := 0
	for _, block := range n.Chain {
		fmt.Println(block.To_string())
		//find a way to get transactions in order from the merkle tree
		if block.MTree == nil {
			fmt.Println("No Transactions In Block Yet")
			return
		}
		for _, l := range block.MTree.Leafs {
			totalTrans += 1
			trans := structures.Deserialize(l.HashedData)
			fmt.Println(trans)
		}
		fmt.Println()
	}

	fmt.Printf("Number of Transactions %d\n", totalTrans)

}

//This function will give the option to display all posts in the chain
//      Also displays author (public key) of the post
//Verify the chain
//Make a post
//Node can be killed with ctrl-d
func (n *Node) Cli_prompt() {

	reader := bufio.NewReader(os.Stdin) //create a reader to parse input

	options := map[string]func(){
		"list":   n.Print_chain,
		"verify": n.Verify_chain,
		"post":   n.Create_transaction,
	}

	//n.Killed is just there in the case we want to kill it from other functions
	for !n.Killed {
		fmt.Println("Enter a command: (list, verify, post)")

		command, err := reader.ReadString('\n')

		if err == io.EOF {
			n.Exit()
			return
		}
		//Clear the newline from input
		command = strings.Replace(command, "\n", "", -1)

		fmt.Println()
		//check if our command is valid
		found := false
		for k, v := range options {
			if command == k {
				found = true
				v()
				break
			}
		}
		if !found {
			fmt.Println("Invalid command, please try again.")
			continue
		}

	}

}

func (n *Node) Print_peer_completions() {
	for _, V := range n.Peer_completions {
		fmt.Printf("\n Peer %d completed the block %d \n", V.Peer, V.BlockIndex)
	}
}



