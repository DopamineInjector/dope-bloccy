package controller

import (
	"dope-bloccy/nft"
	"dope-bloccy/node"
	"dope-bloccy/repository"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetUserNft(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	id := r.PathValue(ID_PATH_PARAM)
	if id == "" {
		http.Error(w, "no id provided", http.StatusNotFound)
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "invalid uuid provided", http.StatusBadRequest)
		return
	}
	user, err := repository.GetUser(id, db)
	if err != nil {
		http.Error(w, "error while fetching data from db", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "no such user", http.StatusNotFound)
		return
	}
	nftIds, err := node.GetUserNfts(user.PubKey)
	if err != nil {
		http.Error(w, "error while getting user nfts", http.StatusInternalServerError)
		return
	}
	metadataIds := node.GetNftMetadataParallel(user.PubKey, nftIds)
	nfts := nft.GetUserNfts(metadataIds, nftIds)
	body, err := json.Marshal(nfts)
	if err != nil {
		http.Error(w, "error encoding response", http.StatusInternalServerError)
		return
	}
	w.Header().Add("content-type", "application/json")
	w.Write(body)
	return
}

func MintNft(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	body := &nft.MintNftRequest{}
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Warn(err.Error())
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(requestBody, body)
	if err != nil {
		log.Warn(err.Error())
		log.Warn(string(requestBody))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	user, err := repository.GetUser(body.User, db)
	if err != nil {
		http.Error(w, "error while fetching data from db", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "no such user", http.StatusNotFound)
		return
	}
	metadata, err := nft.MintNft(body)
	if err != nil {
		http.Error(w, "error while minting nft", http.StatusInternalServerError)
		return
	}
	err = node.MintNft(user.PubKey, metadata.Id)
	if err != nil {
		http.Error(w, "error while minting nft in blockchain node", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func TransferNft(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	body := &nft.TransferNftRequest{}
	requestBody, err := io.ReadAll(r.Body)
	if err != nil {
		log.Warn(err.Error())
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(requestBody, body)
	if err != nil {
		log.Warn(err.Error())
		log.Warn(string(requestBody))
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	sender, err := repository.GetUser(body.Sender, db)
	if err != nil {
		http.Error(w, "no sender wallet data found", http.StatusNotFound)
		return
	}
	recipient, err := repository.GetUser(body.Recipient, db)
	if err != nil {
		http.Error(w, "no recipient wallet data found", http.StatusNotFound)
		return
	}
	err = node.TransferNft(sender.PubKey, recipient.PubKey, body.TokenId, sender.PrivKey);
	if err != nil {
		http.Error(w, "error while transferring nft", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}

