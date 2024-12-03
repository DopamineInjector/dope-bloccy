package node

import (
	"bytes"
	"dope-bloccy/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func getServerAddress() string {
	return utils.GetConfigString(utils.NodeAddress)
}

func GetAccountBalance(walletId []byte) (float32, error) {
	const BALANCE_ENDPOINT = "/api/account/info"
	address := getServerAddress()
	url := fmt.Sprintf("%s/%s", address, BALANCE_ENDPOINT)
	data := GetAccountInfoDto{
		WalletId: walletId,
	}
	requestBody, _ := json.Marshal(data)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(requestBody))
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
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

func TransferFunds(sender, recipient []byte, amount float32, senderPrivKey []byte) error {
	const TRANSFER_ENDPOINT = "/api/transfer"
	address := getServerAddress()
	url := fmt.Sprintf("%s/%s", address, TRANSFER_ENDPOINT)
	payload := TransferRequestPayload{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
	}
	signature, err := signUserTransaction(payload, senderPrivKey)
	if err != nil {
		return err
	}
	body := TransferRequest{
		Payload:   payload,
		Signature: signature,
	}
	requestBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("could not verify signature")
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("some errror that should never occur but yet here we are")
	}
	return nil
}

func MintNft(recipient []byte, metadataId string) error {
	const SC_ENDPOINT = "/api/smartContract"
	address := getServerAddress()
	url := fmt.Sprintf("%s/%s", address, SC_ENDPOINT)
	metadataServerAddress := utils.GetConfigString(utils.MetadataServer)
	const METADATA_ENDPOINT = "metadata"
	metadataUri := fmt.Sprintf("%s/%s/%s", metadataServerAddress, METADATA_ENDPOINT, metadataId)
	args := MintNftArgs{
		MetadataUri: metadataUri,
		Recipient:   recipient,
	}
	jsonArgs, _ := json.Marshal(args)
	scAddress := utils.GetConfigString(utils.NodeNftAddress)
	adminId := utils.GetConfigString(utils.NodePublicKey)
	payload := SmartContractRequestPayload{
		Entrypoint: "_mint",
		Args:       string(jsonArgs),
		Sender:     []byte(adminId),
		Contract:   []byte(scAddress),
	}
	signature := signAdminTransaction(payload)
	body := SmartContractRequest{
		Payload:   payload,
		Signature: signature,
		IsView:    false,
	}
	requestBody, _ := json.Marshal(body)
	log.Debug(requestBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Warn("Could not reach blockchain node")
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden {
		log.Warn("MintNft: signature verification error - should not happen. ever.")
		return fmt.Errorf("could not verify signature")
	} else if resp.StatusCode != http.StatusOK {
		log.Warn("MintNft: error in response from node")
		return fmt.Errorf("some errror that should never occur but yet here we are")
	}
	return nil
}

func GetUserNfts(userId []byte) ([]int, error) {
	const SC_ENDPOINT = "/api/smartContract"
	address := getServerAddress()
	url := fmt.Sprintf("%s/%s", address, SC_ENDPOINT)
	args := OwnedByArgs{
		Owner: userId,
	}
	jsonAgs, _ := json.Marshal(args)
	scAddress := utils.GetConfigString(utils.NodeNftAddress)
	payload := SmartContractRequestPayload{
		Entrypoint: "_owned_by",
		Args:       string(jsonAgs),
		Sender:     userId,
		Contract:   []byte(scAddress),
	}
	body := SmartContractRequest{
		Payload:   payload,
		Signature: nil,
		IsView:    true,
	}
	requestBody, _ := json.Marshal(body)
	log.Debug(requestBody)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Warn("Could not reach blockchain node")
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Warn("GetUserNfts: node has rejected this stuff for some God forsaken reason, may He have mercy on our souls")
		return nil, fmt.Errorf("i don't know what happened")
	}
	responseBody := &SmartContractResponse{}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Warn("For some reason api is giving bad response.")
		return nil, err
	}
	err = json.Unmarshal(bodyBytes, responseBody)
	if err != nil {
		log.Warn("For some reason api is giving bad response. I mean, bad struct")
		return nil, err
	}
	tokens := make([]int, 0)
	err = json.Unmarshal([]byte(responseBody.Output), &tokens)
	if err != nil {
		log.Warn(responseBody.Output)
		log.Warn("Could not parse int slice from this res")
		return nil, err
	}
	return tokens, nil
}
