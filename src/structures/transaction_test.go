package structures

import (
	"testing"
	"log"
)

func Test_sign_transaction(t *testing.T){
	text := "Sample transaction"
	transaction := CreateTransaction(text, 0)
	cipherText := signTransaction(transaction)
	err := readTransaction(transaction, cipherText)

	if err != nil {
		log.Fatal("Transaction not verified")
	}
}