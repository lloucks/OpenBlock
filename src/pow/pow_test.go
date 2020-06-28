

package pow



import (
	"fmt"
	"structures"
	"testing"
)



func Test_low_difficulty(t *testing.T){
    difficulty := 5

    empty_block := structures.Block{}

    hash := solve_block(empty_block, difficulty)

    fmt.Println("hash is ", hash)


}
