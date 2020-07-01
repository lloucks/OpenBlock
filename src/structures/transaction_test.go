package structures

import (
	"testing"
	"log"
	"bytes"
)

func Test_sign_transaction(t *testing.T){
	text := "Sample transaction"
	transaction := CreateTransaction(text, 0)
	cipherText := signTransaction(transaction)

	decryptedText := readTransaction(transaction, cipherText)

	if bytes.Compare([]byte(text), decryptedText) != 0 {
		log.Fatal("Did not fully decrypt.\n", []byte(text), "\n", decryptedText)

	}
}