package utils

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

func LogRequest(w http.ResponseWriter, r *http.Request) {
	log.Infof("%s: %s", r.Method, r.URL)
}
