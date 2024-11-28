package main

import (
	"dope-bloccy/repository"
	"dope-bloccy/utils"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Reading config file")
	err := utils.ReadConfig()
	if err != nil {
		log.Warnf("Could not read config file, reverting to defaults: %s", err.Error())
	}
	db, err := repository.InitializeDBConnection()
	if err != nil {
		log.Fatalf("Could not connect to db: %s", err.Error())
	}
	log.Info("Sucessfully connected to the db")
	router := Routes(db)
	port := utils.GetConfigString(utils.ServerPort)
	log.Infof("Starting to listen on port %s", port)
	port = fmt.Sprintf(":%s", port)
	err = http.ListenAndServe(port, router)
	if err != nil {
		log.Fatalf("Could not start server: %s", err)
	}
}
