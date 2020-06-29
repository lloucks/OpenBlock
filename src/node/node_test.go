

package node



import (
	"fmt"
	"log"
	"structures"
	"testing"
	"time"
)



func Test_difficulty_adjustment_positive(t *testing.T){
    fmt.Printf("Starting difficulty adjustment test A\n")

    //generate 100 blocks 10 seconds apart and add them to the chain
    //We are not going to verify they are valid in this test

    node := Make_node()
    node.Block_time = 20*time.Second
    node.Difficulty = 5

    cur_time := time.Now()

    for i:=0; i<100; i++{
        block := structures.Block{}
        block.Header.Difficulty = 4
        block.Header.Timestamp = cur_time.Add(time.Second * time.Duration(10 * i))
        node.Chain = append(node.Chain, block)
    }

    node.Adjust_difficulty()

    if node.Difficulty != 6{
        log.Fatalf("Difficulty should have been adjusted to 6")
    }

}


func Test_difficulty_adjustment_negative(t *testing.T){
    fmt.Printf("Starting difficulty adjustment test B\n")

    node := Make_node()
    node.Block_time = 20*time.Second
    node.Difficulty = 5

    cur_time := time.Now()

    for i:=0; i<100; i++{
        block := structures.Block{}
        block.Header.Difficulty = 4
        block.Header.Timestamp = cur_time.Add(time.Second * time.Duration(30 * i))
        node.Chain = append(node.Chain, block)
    }

    node.Adjust_difficulty()

    if node.Difficulty != 4{
        log.Fatalf("Difficulty should have been adjusted to 4")
    }



}


//this test will do an actual run and see how long it takes to achieve the proper block time
func Test_difficulty_adjustment_long(t *testing.T){
    fmt.Printf("Starting block generation difficulty adjustement test\n")


}
