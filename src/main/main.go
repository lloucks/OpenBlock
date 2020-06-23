package main

import "fmt"
import "keys"

func main() {
	fmt.Println("Generating RSA key pair and committing to local storage.")
	keys.GenerateKeys()
}
