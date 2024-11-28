package node

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"dope-bloccy/utils"
	"encoding/base64"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

func getAdminPrivateKey() *rsa.PrivateKey {
	stringified := utils.GetConfigString(utils.NodePrivateKey)
	bytes, err := base64.StdEncoding.DecodeString(stringified)
	if err != nil {
		log.Fatalf("could not decode block from privkey in config, %s", err.Error())
	}
	res, err := x509.ParsePKCS1PrivateKey(bytes)
	if err != nil {
		log.Fatalf("Could not parse private key from config file, %s", err.Error())
	}
	return res
}

func signUserTransaction(payload interface{}, privKey []byte) ([]byte, error) {
	key, err := x509.ParsePKCS1PrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	stringified, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	hashed := sha256.Sum256(stringified)
	signed, err := key.Sign(rand.Reader, hashed[:], crypto.SHA256)
	return signed, err
}

func signAdminTransaction(payload interface{}) []byte {
	key := getAdminPrivateKey()
	stringified, err := json.Marshal(payload)
	if err != nil {
		log.Fatalf("Error while marshalling payload of admin transactin: %s", err.Error())
		return nil
	}
	hashed := sha256.Sum256(stringified)
	signed, err := key.Sign(rand.Reader, hashed[:], crypto.SHA256)
	if err != nil {
		log.Fatalf("Error while signing admin transaction payload: %s", err.Error())
	}
	return signed
}
