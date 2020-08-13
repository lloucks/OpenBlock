package keys

import "crypto/x509"
import "crypto/rsa"
import "io/ioutil"
import "crypto/rand"
import "encoding/pem"
import "log"
import "errors"

func GenerateKeys() *rsa.PrivateKey{
	//Use this if we don't want to putput the keys to a file.
	Priv, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil{
		log.Fatal("Failed to generate key pair.")
	}
	return Priv
}

//Generate keys using PKCS#1 ASN.1 specifications
// and write to local file.
func GenerateKeys_toFile() (*rsa.PrivateKey){
	Priv, err := rsa.GenerateKey(rand.Reader, 4096)

	if err != nil {
	   log.Fatal("Failed to generate key pair.")
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

	return Priv
}

//Will either return keys found from local storage,
//or will create new ones.
func GetKeys() (*rsa.PrivateKey, *rsa.PublicKey) {

	err := keysPresent()
	if err != nil{
		GenerateKeys()
	}

	privContent, err := ioutil.ReadFile("key.priv")
	pubContent, err2 := ioutil.ReadFile("key.pub")

	if err != nil || err2 != nil {
		return nil, nil
	}

	privBlock, _ := pem.Decode(privContent)
	pubBlock, _ := pem.Decode(pubContent)
	
	privKey, err := x509.ParsePKCS1PrivateKey(privBlock.Bytes)
	pubKey, err2 := x509.ParsePKCS1PublicKey(pubBlock.Bytes)

	if err != nil || err2 != nil {
		log.Fatal("Could not read keys\n", err, "\n", err2)
	}

	return privKey, pubKey
}

//Check to see if keys are present in the directory.
func keysPresent() error {
	_, err := ioutil.ReadFile("key.priv")
	_, err2 := ioutil.ReadFile("key.pub")

	if err != nil || err2 != nil {
		return errors.New("No keys found")
	}

	return nil

}
