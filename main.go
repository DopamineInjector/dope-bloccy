package main

import (
	"dope-bloccy/repository"
	"dope-bloccy/utils"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Reading config file");
  err := utils.ReadConfig();
  if err != nil {
    log.Warnf("Could not read config file, reverting to defaults: %s", err.Error());
  }
	_, err = repository.InitializeDBConnection();
	if err != nil {
		log.Fatalf("Could not connect to db: %s", err.Error())
	}
	log.Info("Sucessfully connected to the db");
	var forever chan struct{};
	<-forever
}
