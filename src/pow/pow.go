package pow

import (
	"crypto/sha256"
	"fmt"
	"structures"
        "time"
)

//difficulty will be the number of zeroes to match

//we will start the nonce at zero in our client

func to_hex_string(item interface{}) string{

    bytearr := []byte(fmt.Sprintf("%v", item))
    hex := fmt.Sprintf("%x", bytearr)

    return hex
}


func Complete_block(block structures.Block) structures.Block{
    //set timestamp

    block.Header.Timestamp = time.Now()

    //pull off header for the PoW algorithim

    header := &block.Header


    for{
        if Verify_work(*header){
            fmt.Println("Found the proper hash to solve block!:", GenerateHash(*header))
            return block
        }
        //not enough zeroes? increment nonce and try again
        header.Nonce++

    }

}

func GenerateHash(header structures.BlockHeader) string{

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

    return hashstr

}



func Verify_work(header structures.BlockHeader) bool{

    hashstr := GenerateHash(header)


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
