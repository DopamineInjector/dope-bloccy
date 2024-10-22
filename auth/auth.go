package auth

import (
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"dope-bloccy/utils"
	"encoding/base64"
	"encoding/pem"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const AUTH_HEADER = "x-auth-signature"

func VerifyAuthHeader(w http.ResponseWriter, r *http.Request) bool {
  // Overrides entire auth
  if !utils.GetConfigBool(utils.AuthEnabled) {
    log.Warnf("Auth is disabled while handling %s request for %s", r.Method, r.URL)
    return true
  }
  signature := r.Header.Get(AUTH_HEADER);
  if signature == "" {
    http.Error(w, "no auth header present", http.StatusUnauthorized);
    return false;
  }
  var byteSig []byte;
  _, err := base64.StdEncoding.Decode([]byte(signature), byteSig)
  if err != nil {
    http.Error(w, "error while parsing auth signature", http.StatusUnauthorized);
    return false;
  }
  body, err := io.ReadAll(r.Body);
  if err != nil {
    http.Error(w, "error while parsing auth signature", http.StatusUnauthorized);
    return false;
  }
  log.Info(string(body))
  if !verifySignature(byteSig, body) {
    http.Error(w, "could not authorize request", http.StatusForbidden);
    return false;
  }
  return true
}

func verifySignature(signature, payload []byte) bool {
  pemKey := utils.GetConfigString(utils.AuthKey);
  block, _ := pem.Decode([]byte(pemKey));
  parseResult, err := x509.ParsePKIXPublicKey(block.Bytes);
  if err != nil {
    log.Warn(err.Error());
    return false
  }
  publicKey, ok := parseResult.(*rsa.PublicKey);
  if !ok {
    log.Warn("Key is not an rsa key");
    return false
  }
  err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, payload, signature);
  if err != nil {
    log.Warn(err.Error())
  }
  return err == nil
}
