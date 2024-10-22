package controller

import (
	"dope-bloccy/nft"
	"dope-bloccy/repository"
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func GetUserNft(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
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
    http.Error(w, "error while fetching data from db", http.StatusInternalServerError);
    return;
  }
  if user == nil {
    http.Error(w, "no such user", http.StatusNotFound);
    return;
  }
  nfts := nft.GetUserNfts(user.PubKey);
  body, err := json.Marshal(nfts);
  if err != nil {
    http.Error(w, "error encoding response", http.StatusInternalServerError);
    return;
  }
  w.Header().Add("content-type", "application/json")
  w.Write(body);
  return;
}

func MintNft(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
  body := &nft.MintNftRequest{};
  requestBody, err := io.ReadAll(r.Body);
  if err != nil {
    log.Warn(err.Error())
    http.Error(w, "invalid request", http.StatusBadRequest);
    return;
  }
  err = json.Unmarshal(requestBody, body);
  if err != nil {
    log.Warn(err.Error())
    log.Warn(string(requestBody));
    http.Error(w, "invalid request", http.StatusBadRequest);
    return;
  }
  user, err := repository.GetUser(body.User, db);
  if err != nil {
    http.Error(w, "error while fetching data from db", http.StatusInternalServerError);
    return;
  }
  if user == nil {
    http.Error(w, "no such user", http.StatusNotFound);
    return;
  }
  err = nft.MintNft(body);
  if err != nil {
    http.Error(w, "error while minting nft", http.StatusInternalServerError);
    return;
  }
  w.WriteHeader(http.StatusCreated);
  return
}
