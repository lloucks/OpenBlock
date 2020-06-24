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

	// MarshalPKCS1__Key converts a key to PKCS#1, ASN.1 DER form.
	pubBytes := pem.EncodeToMemory( &pem.Block{
	    Type:  "RSA PUBLIC KEY",
	    Bytes: x509.MarshalPKCS1PublicKey(&Priv.PublicKey),
	})

	//WriteFile(filename string, data []byte, perm os.FileMode)
	ioutil.WriteFile("key.pub", pubBytes, 0644)

	privBytes := pem.EncodeToMemory( &pem.Block{
	    Type:  "RSA PRIVATE KEY",
	    Bytes: x509.MarshalPKCS1PrivateKey(Priv),
	})

	ioutil.WriteFile("key.priv", privBytes, 0644)
}
