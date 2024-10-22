package nft

import (
	"dope-bloccy/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const GET_METADATA_ENDPOINT = "metadata"

func getServerAddress() string {
	return utils.GetConfigString(utils.MetadataServer);
}

func getNft(id string) (*NftMetadata, error) {
  address := getServerAddress();
  url := fmt.Sprintf("%s/%s/%s", address, GET_METADATA_ENDPOINT, id);
  log.Info(url)
  resp, err := http.Get(url);
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close();
  body, err := io.ReadAll(resp.Body);
  if err != nil {
    return nil, err
  }
  res := &NftMetadata{};
  err = json.Unmarshal(body, res);
  return res, err
}
