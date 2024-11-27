package repository

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"

	"gorm.io/gorm"
)

func AddUser(id string, db *gorm.DB) error {
  pubkey, privkey, err := generateUserKeys();
  if err != nil {
    return err
  }
  err = db.Create(&User{
    ID: id,
    PubKey: pubkey,
    PrivKey: privkey,
  }).Error
  if err != nil {
    return &UserExistsError{}
  }
  return nil
}

func GetUser(id string, db *gorm.DB) (user *User, err error) {
  user = &User{}
  err = db.First(user, "id=?", id).Error;
  if err != nil {
    println(err.Error())
  }
  return
}

func generateUserKeys() (pubkey, privkey []byte, err error) {
  priv, err := rsa.GenerateKey(rand.Reader, 2048);
  if err != nil {
    return
  }
  pub := priv.PublicKey;
  privkey = x509.MarshalPKCS1PrivateKey(priv);
  pubkey = x509.MarshalPKCS1PublicKey(&pub);
  return
}
