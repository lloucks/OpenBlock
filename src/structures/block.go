//this pakage will define blocks and transactions, plus any methods we need to work with them
//such as a method to create a block by passing in transactions, appending transactions etc...


package structures
//highly based off of https://github.com/bitcoin/bitcoin/tree/master/src/primitives

import (
    "time"
)

//These data types are not representative of the actual product

type Block struct{
    Index int
    //Data string
    //Data is stored in transactions, not blocks.
    Magic_num uint32 //so we know it's part of our protocol
    Block_size uint32 //number of bytes in this block
    Header BlockHeader
    T_count uint32 //number of transactions in this block
    Transactions []Transaction
}

type BlockHeader struct{
    Prev_block_hash int
    //Hash int
    //you cannot do this
    Merkle_root_hash int
    Timestamp time.Time //created by the node making this block
    Difficulty uint32 //the difficulty, how many 0s need to be in PoW hash
    Nonce uint32 //the number that solves this block
}


//These functions are already defined in pow.go
/*

func CreateBlock() Block{
    var block Block

    //block.Magic_num = 0x8b665966a5e97c796e268734cfe5eaf0
    block.Magic_num = 4 //why the hash?? this is just used as an identifier.
    //block.Block_size = 0
    block.Header.Timestamp = time.Now() //aready defined

    return block
}
func getPrefix(length int) string {
    letterBytes := 0

    b:= make([]byte, length)

    for i := range b{
        b[i] = letterBytes[0]
    }
    return string(b)
}



func (b *Block) Mine(){
    prefix := getPrefix(b.Difficulty)

    for{
        b.GenerateHash()

        if strings.HasPrefix(b.Header.Hash.String(),prefix){
            break
        } else {
            b.Header.Nonce = b.Header.Nonce + 1
            b.GenerateHash()
        }
    }
}
*/
