

package pow



import (
	"fmt"
	"log"
	"structures"
	"testing"
)



func Test_low_difficulty(t *testing.T){
    fmt.Println("Starting low difficulty test(2-20 seconds)")
    empty_block := structures.Block{}

    empty_block.Header.Difficulty = 5

    block := complete_block(empty_block)

    fmt.Println("Nonce is ", block.Header.Nonce)

    if verify_work(block.Header){
        fmt.Printf("Block is valid\n")
    }else{
        log.Fatalf("Block is not valid")
    }
}
