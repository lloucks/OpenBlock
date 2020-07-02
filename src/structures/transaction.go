package structures

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/gob"
	"keys"
	"log"
)

//Given that we are not using currency, the transaction strucutre is flexible.

//We can make it so that we have some designated address that all posts are sent to

//Please add to this structure when required to make a part work
type Transaction struct {
	ID     []byte
	Text   string //contents of the post
	Author int    //Will be the authors signature, that other nodes should be able to verify
	//int is just a placeholder for now

	privateKey rsa.PrivateKey
	publicKey  rsa.PublicKey
	//need to send it somewhere
	//need to sign it so it can be veried
}

func signTransaction(transaction *Transaction) []byte {
	privKey, pubKey := keys.GetKeys()

	transaction.privateKey = *privKey
	transaction.publicKey = *pubKey

	transactionBytes := transaction.Serialize()
	hashed := sha256.Sum256(transactionBytes)
	signature, _ := rsa.SignPKCS1v15(rand.Reader, &transaction.privateKey, crypto.SHA256, hashed[:])

	return signature
}

func readTransaction(transaction *Transaction, signature []byte) error {
	hashed := sha256.Sum256(transaction.Serialize())
	err := rsa.VerifyPKCS1v15(&transaction.publicKey, crypto.SHA256, hashed[:], signature)

	return err
}

func CreateTransaction(text string, author int) *Transaction {
	t := Transaction{}
	t.Text = text
	t.Author = author
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
