package pow

import (
	"crypto/sha256"
	"fmt"
	"structures"
	"time"
	//        "encoding/binary"
)

//difficulty will be the number of zeroes to match

//we will start the nonce at zero in our client

func to_hex_string(item interface{}) string {

	bytearr := []byte(fmt.Sprintf("%v", item))
	hex := fmt.Sprintf("%x", bytearr)

	return hex
}

func Complete_block(block structures.Block) structures.Block {
	//set timestamp

	block.Header.Timestamp = time.Now()

	//pull off header for the PoW algorithim

	header := &block.Header

	for {
		if Verify_work(*header) {
			//fmt.Println("Validated hash: ", to_bit_string(GenerateHash(*header)))
			return block
		}
		//not enough zeroes? increment nonce and try again
		header.Nonce++

	}

}

func GenerateHash(header structures.BlockHeader) [32]byte {

	prevhex := to_hex_string(header.Prev_block_hash)
	merklehex := to_hex_string(header.Merkle_root_hash)
	timehex := to_hex_string(header.Timestamp)
	bitshex := to_hex_string(header.Difficulty)
	noncehex := to_hex_string(header.Nonce)

	headerhex := prevhex + merklehex + timehex + bitshex + noncehex

	bytes := []byte(fmt.Sprintf("%v", headerhex)) //

	hash := sha256.Sum256(bytes)

	return hash

}

func Verify_work(header structures.BlockHeader) bool {

	//should not all 0 Difficulty because then nodes can send empty blocks and they get veirified as valid
	if header.Difficulty == 0 {
		return false
	}
	hash := GenerateHash(header)

	hashstr := to_bit_string(hash)

	difficulty := int(header.Difficulty)

	cmpstring := ""

	for i := 0; i < difficulty; i++ {
		cmpstring = cmpstring + "0"
	}
	if hashstr[:difficulty] == cmpstring {
		return true
	} else {
		return false
	}

}

func to_bit_string(hash [32]byte) string {
	output := ""

	for _, b := range hash {
		output = output + fmt.Sprintf("%08b", b)
	}

	return output
}
