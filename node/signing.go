package node

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"dope-bloccy/utils"
	"encoding/json"

	log "github.com/sirupsen/logrus"
)

func getAdminPrivateKey() *rsa.PrivateKey {
  stringified := utils.GetConfigString(utils.NodePrivateKey);
  res, err := x509.ParsePKCS1PrivateKey([]byte(stringified));
  if err != nil {
    log.Fatal("Could not parse private key from config file") 
  }
  return res
}

func SignUserTransaction(payload interface{}, privKey string) ([]byte, error) {
  key, err := x509.ParsePKCS1PrivateKey([]byte(privKey));
  if err != nil {
    return nil, err
  }
  stringified, err := json.Marshal(payload);
  if err != nil {
    return nil, err
  }
  signed, err := key.Sign(rand.Reader, stringified, crypto.SHA256);
  return signed, err
}

func SignAdminTransaction(payload interface{}) ([]byte, error) {
  key := getAdminPrivateKey();
  stringified, err := json.Marshal(payload);
  if err != nil {
    return nil, err
  }
  signed, err := key.Sign(rand.Reader, stringified, crypto.SHA256);
  return signed, err
}
