package structures

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"keys"
	"log"
	"errors"
)

//Given that we are not using currency, the transaction strucutre is flexible.
//We can make it so that we have some designated address that all posts are sent to
//Please add to this structure when required to make a part work
type Transaction struct {
	ID     []byte
	Text   string //contents of the post
	Author int    //Will be the authors signature, that other nodes should be able to verify
	//int is just a placeholder for now
	Data []byte

	Signature []byte
	publicKey  rsa.PublicKey
	//need to send it somewhere
}

func SignTransaction_withoutFile(transaction *Transaction, privKey *rsa.PrivateKey) []byte {
	transactionBytes := transaction.Serialize()
	hashed := sha256.Sum256(transactionBytes)
	signature, _ := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed[:])

	return signature
}

//returns a signature, which is the signed serialization of the transaction
func SignTransaction(transaction *Transaction) []byte {
	privKey, pubKey := keys.GetKeys()
	transaction.publicKey = *pubKey

	//transactionBytes := transaction.Serialize()
	hashed := sha256.Sum256(transaction.Data)
	signature, _ := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed[:])
	return signature
}

func VerifyTransaction_withoutFile(transaction *Transaction) error {
	if transaction.publicKey.N == nil {
		return errors.New("Public key is empty.")
	}
	hashed := sha256.Sum256(transaction.Serialize())
	err := rsa.VerifyPKCS1v15(&transaction.publicKey, crypto.SHA256, hashed[:], transaction.Signature)

	return err
}

//verifies a transaction given the signature
func VerifyTransaction(transaction *Transaction) error {
	//replace with RPC call when distributed.
	_, pubKey := keys.GetKeys()

	hashed := sha256.Sum256(transaction.Data)
	err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], transaction.Signature)
	return err
}

func (t Transaction) UpdateTransactionHash() *Transaction {
 	hash := t.Serialize()
 	t.ID = hash[:]
 	return &t
 }

func CreateTransaction(text string, author int) *Transaction {
	t := Transaction{}
	t.Text = text
	t.Author = author
	t.Data = t.Serialize()
	t.Signature = SignTransaction(&t)
	hash := t.Serialize()
	t.ID = hash[:]
	return &t
}

func (t Transaction) Serialize() []byte {
	buf := &bytes.Buffer{}
	if err := gob.NewEncoder(buf).Encode(t); err != nil {
		log.Panic(err)
		return nil
	}
	return buf.Bytes()
}

func Deserialize(serialized []byte) *Transaction {
	var result Transaction
	d := gob.NewDecoder(bytes.NewReader(serialized))
	err := d.Decode(&result)
	if err != nil {
		log.Fatal("Error Deserializing Transaction")
		return nil
	}
	return &result
}

func (t *Transaction) to_string() string {
    //Do I need to deserialize first???
    var result string

    result += fmt.Sprintf("Author: %v\n", t.publicKey)
    result += fmt.Sprintf("Post:\n      %v\n", t.publicKey)

    return result
}
