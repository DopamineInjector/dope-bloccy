package main

import (
	"dope-bloccy/auth"
	"dope-bloccy/controller"
	"dope-bloccy/utils"
	"log"
	"net/http"

	"gorm.io/gorm"
)

func Routes(db *gorm.DB) (router *http.ServeMux) {
	if db == nil {
		log.Panic("db connection is nil while creating routes")
	}
	router = http.DefaultServeMux
	router.HandleFunc("POST /api/wallet/{id}", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(w, r)
		if !auth.VerifyAuthHeader(w, r) {
			return
		}
		controller.HandleAddUser(w, r, db)
	})
	router.HandleFunc("GET /api/wallet/{id}", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(w, r)
		controller.HandleGetUser(w, r, db)
	})
	router.HandleFunc("POST /api/transfer", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(w, r)
		controller.HandleTransfer(w, r, db)
	})
	// Nft endpoints
	router.HandleFunc("GET /api/wallet/{id}/nfts", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(w, r)
		controller.GetUserNft(w, r, db)
	})
	router.HandleFunc("POST /api/nft/mint", func(w http.ResponseWriter, r *http.Request) {
		utils.LogRequest(w, r)
		if !auth.VerifyAuthHeader(w, r) {
			return
		}
		controller.MintNft(w, r, db)
	})
	return
}
