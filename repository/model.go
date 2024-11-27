package repository

import (
	"time"
)

type User struct {
  ID string
  CreatedAt time.Time
  UpdatedAt time.Time
  PubKey []byte
  PrivKey []byte
}
