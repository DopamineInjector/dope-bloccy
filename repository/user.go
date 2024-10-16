package repository

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"

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

func generateUserKeys() (pubkey, privkey string, err error) {
  priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader);
  if err != nil {
    return
  }
  pub := priv.PublicKey;
  privkey = hex.EncodeToString(priv.D.Bytes());
  pubkey = hex.EncodeToString(append(pub.X.Bytes(), pub.Y.Bytes()...));
  return
}
