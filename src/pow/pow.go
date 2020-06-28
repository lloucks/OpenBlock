//this pakage will define blocks and transactions, plus any methods we need to work with them



package pow
//highly based off of https://github.com/bitcoin/bitcoin/tree/master/src/primitives

import (
	"crypto/sha256"
	"fmt"
	"structures"
        "time"
)

//These data types are not representative of the actual product

//difficulty will be the number of zeroes to match

//we will start the nonce at zero in our client

func to_hex_string(item interface{}) string{

    bytearr := []byte(fmt.Sprintf("%v", item))
    hex := fmt.Sprintf("%x", bytearr)

    return hex
}


func complete_block(block structures.Block) structures.Block{
    //set timestamp

    block.Header.Timestamp = time.Now()

    //pull off header for the PoW algorithim

    header := &block.Header


    for{
        if verify_work(*header){
            return block
        }
        //not enough zeroes? increment nonce and try again
        header.Nonce++

    }

}


func verify_work(header structures.BlockHeader) bool{

    prevhex := to_hex_string(header.Prev_block_hash)
    merklehex := to_hex_string(header.Merkle_root_hash)
    timehex := to_hex_string(header.Timestamp)
    bitshex := to_hex_string(header.Difficulty)
    noncehex := to_hex_string(header.Nonce)

    headerhex := prevhex + merklehex + timehex + bitshex + noncehex

    bytes := []byte(fmt.Sprintf("%v", headerhex)) //

    //make the hash

    hash := sha256.Sum256(bytes)

    //convert back to hex

    hashstr := fmt.Sprintf("%x", hash)

    //good number of zeroes?

//    fmt.Println(hashstr)

    difficulty := int(header.Difficulty)

    var cmpstring string //create the string to verify we have correct nonce
    for i := 0; i<int(difficulty); i++{
        cmpstring = cmpstring + "0"
    }

    if hashstr[:difficulty] == cmpstring{
        return true
    } else{
        return false
    }

}
