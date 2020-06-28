//this pakage will define blocks and transactions, plus any methods we need to work with them
//such as a method to create a block by passing in transactions, appending transactions etc...


package pow
//highly based off of https://github.com/bitcoin/bitcoin/tree/master/src/primitives

import (
	"crypto/sha256"
	"fmt"
	"structures"
)

//These data types are not representative of the actual product

//difficulty will be the number of zeroes to match

//we will start the nonce at zero in our client


func solve_block(block structures.Block, difficulty int) [32]byte{
    //convert to byte array


    var cmpstring string

    for x := 0; x<difficulty; x++{
        cmpstring = cmpstring + "0"
    }
    fmt.Println(cmpstring)

    var winner [32]byte
    //hash it
    for{

        bytes := []byte(fmt.Sprintf("%v", block))

        hash := sha256.Sum256(bytes)

        hashstr := fmt.Sprintf("%x", hash)

        //good number of zeroes?
        fmt.Println(hashstr)

        if hashstr[:difficulty] == cmpstring{
            fmt.Println("found a valid hash!")
            winner = hash
            break
        }


        //not enough zeroes? increment nonce and try again
        block.Header.Nonce ++

        fmt.Println("Setting nonce to ", block.Header.Nonce)

    }

    return winner

}


func verify_block(block structures.Block){




}
