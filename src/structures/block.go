//this pakage will define blocks and transactions, plus any methods we need to work with them
//such as a method to create a block by passing in transactions, appending transactions etc...

package structures

//highly based off of https://github.com/bitcoin/bitcoin/tree/master/src/primitives

import (
	"time"
)

type Block struct {
	Index int
	Magic_num  uint32 //so we know it's part of our protocol
	Block_size uint32 //number of bytes in this block
	Header     BlockHeader
	T_count    uint32 //number of transactions in this block
	MTree      *MerkleTree
}

type BlockHeader struct {
	Prev_block_hash int
	Merkle_root_hash int
	Timestamp        time.Time //created by the node making this block
	Difficulty       uint32    //the difficulty, how many 0s need to be in PoW hash
	Nonce            uint32    //the number that solves this block
}
