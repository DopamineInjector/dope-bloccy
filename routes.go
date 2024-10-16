package main

import (
	"dope-bloccy/controller"
	"log"
	"net/http"

	"gorm.io/gorm"
)

func Routes(db *gorm.DB) (router *http.ServeMux) {
  if db == nil {
    log.Panic("db connection is nil while creating routes")
  }
  router = http.DefaultServeMux;
  router.HandleFunc("POST /api/wallet/{id}", func(w http.ResponseWriter, r *http.Request) {
    controller.HandleAddUser(w, r, db);
  })
  return
}
