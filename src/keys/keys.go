package keys

import "crypto/x509"
import "crypto/rsa"
import "io/ioutil"
import "crypto/rand"
import "encoding/pem"
import "fmt"

func GenerateKeys(){
	Priv, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
	   fmt.Println("Failed to generate key pair.")
	}

	// MarshalPKIXPublicKey converts a public key to PKIX, ASN.1 DER form.
	pubASN1, err := x509.MarshalPKIXPublicKey(&Priv.PublicKey)

	if err != nil {
	    fmt.Println("Failed to convert public key.")
	}

	pubBytes := pem.EncodeToMemory( &pem.Block{
	    Type:  "RSA PUBLIC KEY",
	    Bytes: pubASN1,
	})

	ioutil.WriteFile("key.pub", pubBytes, 0644)
}
