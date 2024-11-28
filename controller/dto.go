package controller

import (
	"dope-bloccy/repository"
	"encoding/json"
	"time"
)

// Get wallet info

type GetWalletResponse struct {
	Id        string    `json:"id"`
	PubKey    []byte    `json:"publicKey"`
	Balance   float32   `json:"balance"`
	CreatedAt time.Time `json:"created"`
}

func (t *GetWalletResponse) Json() []byte {
	res, _ := json.Marshal(t)
	return res
}

func ResponseFromUser(user *repository.User, balance float32) *GetWalletResponse {
	return &GetWalletResponse{
		Id:        user.ID,
		PubKey:    user.PubKey,
		CreatedAt: user.CreatedAt,
		Balance:   balance,
	}
}

// Transfer funds

type TransferFundsRequest struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float32 `json:"amount"`
}

func (t *TransferFundsRequest) Json() []byte {
	res, _ := json.Marshal(t)
	return res
}
