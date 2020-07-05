package node

import (
	"fmt"
	"log"
	"pow"
	"structures"
	"testing"
	"time"
)

func Test_difficulty_adjustment_positive(t *testing.T) {
	fmt.Printf("Starting difficulty adjustment test A\n")

	//generate 100 blocks 10 seconds apart and add them to the chain
	//We are not going to verify they are valid in this test

	node := Make_node()
	node.Block_time = 20 * time.Second
	node.Cur_difficulty = 5

	cur_time := time.Now()

	for i := 0; i < 100; i++ {
		block := structures.Block{}
		block.Header.Difficulty = 4
		block.Header.Timestamp = cur_time.Add(time.Second * time.Duration(10*i))
		node.Chain = append(node.Chain, block)
	}

	node.Adjust_difficulty()

	if node.Cur_difficulty != 6 {
		log.Fatalf("Cur_difficulty should have been adjusted to 6")
	}

}

func Test_difficulty_adjustment_negative(t *testing.T) {
	fmt.Printf("Starting difficulty adjustment test B\n")

	node := Make_node()
	node.Block_time = 20 * time.Second
	node.Cur_difficulty = 5

	cur_time := time.Now()

	for i := 0; i < 100; i++ {
		block := structures.Block{}
		block.Header.Difficulty = 4
		block.Header.Timestamp = cur_time.Add(time.Second * time.Duration(30*i))
		node.Chain = append(node.Chain, block)
	}

	node.Adjust_difficulty()

	if node.Cur_difficulty != 4 {
		log.Fatalf("Cur_difficulty should have been adjusted to 4")
	}

}

//this test will do an actual run and see how long it takes to achieve the proper block time
func Test_difficulty_adjustment_long(t *testing.T){
    fmt.Printf("Starting block generation difficulty adjustement test\n")

    node := Make_node()
    node.Block_time = 4 * time.Second
    node.Cur_difficulty = 1 //start it low and wait for it to adjust up past the target

    cur_time := time.Now()

    timeout := cur_time.Add(time.Second * 120) //Should do it in less time than this

    complete := false

    //Generate blocks and complete them until block time is close to 5 Seconds

    for timeout.Sub(time.Now()) > 0{
        block := structures.Block{}
        block.Header.Difficulty = uint32(node.Cur_difficulty)
        node.Chain = append(node.Chain, pow.Complete_block(block))

        if len(node.Chain) % 5 == 0{
            node.Adjust_difficulty()


    //check if last 3 blocks are past the goal

        if len(node.Chain) > 5{
            blocks := node.Chain[len(node.Chain)-6:]
            var times []time.Time

            for _, block := range(blocks){
                times = append(times, block.Header.Timestamp)
            }

            var differences []time.Duration

            for i, t := range(times){
                if i == 0{
                    continue
                }
                differences = append(differences, t.Sub(times[i-1]))
            }

            average := time.Duration(0)
            for _, t := range(differences){
                average += t
            }

            average = average/time.Duration(len(differences))

            fmt.Printf("Average time is %v\n", average)

            if average > time.Duration(time.Second) * 2{ //close enough to goal
                fmt.Println("Complete")
                complete = true
                break
            }
        }


    }
}


    if complete{
        fmt.Println("Difficulty was adjusted enough to get proper block times, PASS")
    } else {
        log.Fatalf("Difficulty was not adjusted well enough, FAIL")
    }


}
