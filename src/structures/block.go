//this pakage will define blocks and transactions, plus any methods we need to work with them
//such as a method to create a block by passing in transactions, appending transactions etc...

package structures

//highly based off of https://github.com/bitcoin/bitcoin/tree/master/src/primitives

import (
	"fmt"
	"time"
)

type Block struct {
	Index      int
	Block_size uint32 //number of bytes in this block
	Header     BlockHeader
	TList      []*Transaction
}

type BlockHeader struct {
	Prev_block_hash  [32]byte
	Merkle_root_hash int
	Timestamp        time.Time //created by the node making this block
	Difficulty       uint32    //the difficulty, how many 0s need to be in PoW hash
	Nonce            uint32    //the number that solves this block
}

func (b *Block) To_string() string {
	var result string

	result += fmt.Sprintf("Block index: %v\n", b.Index)
	result += fmt.Sprintf("Size: %v posts/block\n", b.Block_size)
	result += fmt.Sprintf("Header:\n")
	result += fmt.Sprintf("     Previous block hash: %x\n", b.Header.Prev_block_hash)
	result += fmt.Sprintf("     Time produced: %v\n", b.Header.Timestamp)
	result += fmt.Sprintf("     Merkle root: %x\n", b.Header.Merkle_root_hash)
	result += fmt.Sprintf("     Block difficulty: %v\n", b.Header.Difficulty)
	result += fmt.Sprintf("     Nonce: %v\n\n", b.Header.Nonce)

	return result
}
