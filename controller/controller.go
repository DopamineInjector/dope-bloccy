package controller

import (
	"dope-bloccy/node"
	"dope-bloccy/repository"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

const ID_PATH_PARAM = "id";

func HandleAddUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
  if db == nil {
    http.Error(w, "error while connecting to db", http.StatusInternalServerError);
    return;
  }
  id := r.PathValue(ID_PATH_PARAM);
  if id == "" {
    http.Error(w, "no id provided", http.StatusNotFound);
    return;
  }
  if _, err := uuid.Parse(id); err != nil {
    http.Error(w, "invalid uuid provided", http.StatusBadRequest);
    return;
  }
  err := repository.AddUser(id, db);
  if err != nil {
    http.Error(w, "user already exists", http.StatusConflict);
    return;
  }
  w.WriteHeader(http.StatusCreated);
  return
}


func HandleGetUser(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
  if db == nil {
    http.Error(w, "error while connecting to db", http.StatusInternalServerError);
    return;
  }
  id := r.PathValue(ID_PATH_PARAM);
  if id == "" {
    http.Error(w, "no id provided", http.StatusNotFound);
    return;
  }
  if _, err := uuid.Parse(id); err != nil {
    http.Error(w, "invalid uuid provided", http.StatusBadRequest);
    return;
  }
  user, err := repository.GetUser(id, db);
  if err != nil {
    http.Error(w, "error while fetching from db", http.StatusInternalServerError);
    return;
  }
	balance, err := node.GetAccountBalance(user.PubKey);
  if err != nil {
		logrus.Warn("could not fetch account balance from blockchain node");
		balance = -1;
  }
  body := ResponseFromUser(user, balance).Json();
  w.Header().Add("content-type", "application/json");
  w.WriteHeader(http.StatusOK);
  w.Write(body);
  return
}
