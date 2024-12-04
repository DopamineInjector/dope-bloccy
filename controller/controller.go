package controller

import (
	"dope-bloccy/node"
	"dope-bloccy/repository"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const ID_PATH_PARAM = "id"

func HandleAddUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if db == nil {
		http.Error(w, "error while connecting to db", http.StatusInternalServerError)
		return
	}
	id := r.PathValue(ID_PATH_PARAM)
	if id == "" {
		http.Error(w, "no id provided", http.StatusNotFound)
		return
	}
	if _, err := uuid.Parse(id); err != nil {
		http.Error(w, "invalid uuid provided", http.StatusBadRequest)
		return
	}
	pubkey, err := repository.AddUser(id, db);
	if err != nil {
		http.Error(w, "user already exists", http.StatusConflict)
		return
	}
	err = node.CreateAccount(pubkey);
	if err != nil {
		log.Warnln(err.Error())
		http.Error(w, "error while creating user on da blockchain", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func HandleGetUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if db == nil {
		http.Error(w, "error while connecting to db", http.StatusInternalServerError)
		return
	}
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
		http.Error(w, "error while fetching from db", http.StatusInternalServerError)
		return
	}
	balance, err := node.GetAccountBalance(user.PubKey)
	if err != nil {
		log.Warn("could not fetch account balance from blockchain node")
		balance = -1
	}
	body := ResponseFromUser(user, balance).Json()
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
	return
}

func HandleTransfer(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	if db == nil {
		http.Error(w, "error while connecting to db", http.StatusInternalServerError)
		return
	}
	body := &TransferFundsRequest{}
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
	err = node.TransferFunds(sender.PubKey, recipient.PubKey, body.Amount, sender.PrivKey)
	if err != nil {
		log.Warn(err.Error());
		http.Error(w, "error while transferring funds", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
