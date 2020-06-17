//this pakage will define blocks and transactions, plus any methods we need to work with them
//such as a method to create a block by passing in transactions, appending transactions etc...


package block
//highly based off of https://github.com/bitcoin/bitcoin/tree/master/src/primitives

import (
    "time"
)

//These data types are not representative of the actual product
//I need to implement 256 bit ints to store hashes as well (or import a library)

type Block struct{
    Magic_num uint32 //so we know it's part of our protocol
    Block_size uint32 //number of bytes in this block
    Header BlockHeader
    T_count uint32 //number of transactions in this block
    Transactions []Transaction
}

type BlockHeader struct{
    Prev_block_hash int //should be uint256
    Merkle_root_hash int //should be uint256
    timestamp time.Time //created by the node making this block
    bits uint32 //the difficulty, how many 0s need to be in PoW hash
    Nonce uint32 //the number that solves this block
}




//Given that we are not using currency, the transaction strucutre is flexible.

//We can make it so that we have some designated address that all posts are sent to

//Please add to this structure when required to make a part work
type Transaction struct{
    text string //contents of the post
    author int //Will be the authors signature, that other nodes should be able to verify
    //int is just a placeholder for now


    //need to send it somewhere
    //need to sign it so it can be veried



}
