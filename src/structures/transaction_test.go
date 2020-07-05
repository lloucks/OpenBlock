package structures

import (
	"testing"
	"log"
)

func Test_and_sign_transaction(t *testing.T){
	text := "Sample transaction"
	transaction := CreateTransaction(text, 0)
	signature := signTransaction(transaction)
	err := verifyTransaction(transaction, signature)

	//verified transaction returns nil
	if err != nil {
		log.Fatal("Transaction not verified.")
	}

	if transaction.Text != text {
		log.Fatal("Text was modified.")
	}
}

func Test_invalid_transaction(t *testing.T){
	text := "My Sample transaction"
	transaction := CreateTransaction(text, 0)
	signature := transaction.Serialize()

	//give the verfiy some unsigned data.
	err := verifyTransaction(transaction, signature)

	//should return an error
	if err == nil {
		log.Fatal("Invalid transaction was verified.")
	}
}

