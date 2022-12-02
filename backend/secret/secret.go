package secret

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"log"
	"os"
)

var (
	key *rsa.PrivateKey
)

// init sets rsa key.
func init() {
	key = getKey()
}

// getKeyReader gets io.Reader which contains the RSA key. If not exist it will create it.
func getKey() *rsa.PrivateKey {
	for {
		buf, err := os.ReadFile("key.der")

		// Key file found.
		if err == nil {
			key, err := x509.ParsePKCS1PrivateKey(buf)
			if err != nil {
				log.Fatalln("Failed to parse key.der.")
			}

			return key
		}

		// Key missing or no persmission.
		if os.IsNotExist(err) {
			generateKey()
		} else {
			log.Fatalln("Can not open key file:", err)
		}
	}
}

// generateKey genereates and saves key to key.der.
func generateKey() {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalln("Failed to generate random RSA key:", err)
	}

	buf := x509.MarshalPKCS1PrivateKey(key)

	if err = os.WriteFile("key.der", buf, 0644); err != nil {
		log.Fatalln("Failed to create key file:", err)
	}
}

// SigningMethod returns signing method.
func SigningMethod() string {
	return "RS256"
}

// Private returns the private key.
func PrivateKey() crypto.PrivateKey {
	return key
}

// Public returns the public key corresponding to private.
func PublicKey() crypto.PublicKey {
	return key.Public()
}
