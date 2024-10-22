package controller

import (
	"dope-bloccy/nft"
	"dope-bloccy/repository"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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
  // TODO - obtain user nfts from blockchain node.
  userNfts := [2]string{"6a389c77-6bab-4e0c-98b6-4c25ca23c60e", "2bc027ad-573a-4a58-a049-9d6c04800e5c"};
  nfts := nft.GetNftsMetadata(userNfts[:]);
  body, err := json.Marshal(nfts);
  if err != nil {
    http.Error(w, "error encoding response", http.StatusInternalServerError);
    return;
  }
  w.Header().Add("content-type", "application/json")
  w.Write(body);
  return;
}
