package node

import (
	"bytes"
	"dope-bloccy/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getServerAddress() string {
	return utils.GetConfigString(utils.NodeAddress);
}

func GetAccountBalance(walletId string) (float32, error) {
  const BALANCE_ENDPOINT = "/api/account/info"
	address := getServerAddress()
	url := fmt.Sprintf("%s/%s", address, BALANCE_ENDPOINT)
	data := GetAccountInfoDto{
    WalletId: walletId,
	}
	requestBody, _ := json.Marshal(data)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
  req.Header.Set("content-type", "application/json");
  resp, err := http.DefaultClient.Do(req);
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Error while minting metadata on nft server")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	accInfo := &AccountInfoResponseDto{}
	err = json.Unmarshal(body, accInfo)
	return accInfo.Balance, err
}
