// Almighty Father, You are the source of all wisdom, talents, and skills.
// Thank You for these beautiful gifts, opportunities to serve you, and the work that You have entrusted to me.
// Please open my heart and enlighten my mind, so that I may be fully in tune with Your divine purpose in calling me to the engineering profession.
// Lord God, You are the greatest engineer. Please infuse me even with just the tiniest spark of Your Divine Wisdom so that as I  do my work, it is really your work that is done.
// Loving God, make my heart Your Heart, make my mind Your Mind, and make my hands Your Hands.
// Make me your instrument, so that I may be always conscious and mindful of the fact that on my work depend the lives and properties of my fellow human. Bless me also with the gift of love and sensitivity to respect the people who build and use the products of my work.
// This I ask in Jesus’ name.
// AMEN.
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
