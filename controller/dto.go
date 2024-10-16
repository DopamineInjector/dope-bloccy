package controller

import (
	"dope-bloccy/repository"
	"encoding/json"
	"time"
)

type GetWalletResponse struct {
  Id      string      `json:"id"`
  PubKey  string      `json:"publicKey"`
  CreatedAt time.Time `json:"created"`
}

func (t *GetWalletResponse) Json() []byte {
  res, _ := json.Marshal(t);
  return res
}

func ResponseFromUser(user *repository.User) *GetWalletResponse {
  return &GetWalletResponse{
    Id: user.ID,
    PubKey: user.PubKey,
    CreatedAt: user.CreatedAt,
  }
}
