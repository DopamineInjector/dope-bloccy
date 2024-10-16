package repository

import (
	"dope-bloccy/utils"
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitializeDBConnection() (*gorm.DB, error) {
  host := utils.GetConfigString(utils.PostgresHost);
  port := utils.GetConfigString(utils.PostgresPort);
  user := utils.GetConfigString(utils.PostgresUser);
  password := utils.GetConfigString(utils.PostgresPassword);
  db := utils.GetConfigString(utils.PostgresDb);
  log.Infof("Attempting to connect to %s @ %s", db, host);
  dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Warsaw", host, user, password, db, port);
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{});
	if err != nil {
		return nil, err
	}
	err = conn.AutoMigrate(&User{});
	if err != nil {
		return nil, err
	}
	return conn, nil
}
