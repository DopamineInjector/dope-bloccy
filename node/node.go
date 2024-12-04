package node

import (
	"bytes"
	"dope-bloccy/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
)

func getServerAddress() string {
	return utils.GetConfigString(utils.NodeAddress)
}

func CreateAccount(walletId []byte) error {
	const CREATE_ACCOUNT_ENDPOINT = "api/account"
	address := getServerAddress()
	url := fmt.Sprintf("%s/%s", address, CREATE_ACCOUNT_ENDPOINT)
	log.Info(url)
	body := CreateAccountRequest{
		PublicKey: walletId,
	}
	requestBody, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	req.Header.Set("content-type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Warnf("Node threw up error: %d", resp.StatusCode)
		return fmt.Errorf("some errror that should never occur but yet here we are")
	}
	return nil
}

func GetAccountBalance(walletId []byte) (float32, error) {
	const BALANCE_ENDPOINT = "api/account/info"
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
	const TRANSFER_ENDPOINT = "api/transfer"
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
	} else if resp.StatusCode != http.StatusCreated {
		log.Warn(resp.StatusCode)
		return fmt.Errorf("some errror that should never occur but yet here we are")
	}
	return nil
}

func MintNft(recipient []byte, metadataId string) error {
	const SC_ENDPOINT = "api/smartContract"
	address := getServerAddress()
	url := fmt.Sprintf("%s/%s", address, SC_ENDPOINT)
	args := MintNftArgs{
		MetadataUri: metadataId,
		Recipient:   recipient,
	}
	jsonArgs, _ := json.Marshal(args)
	scAddress := utils.GetConfigString(utils.NodeNftAddress)
	adminId := utils.GetConfigBytes(utils.NodePublicKey)
	payload := SmartContractRequestPayload{
		Entrypoint: "_mint",
		Args:       string(jsonArgs),
		Sender:     adminId,
		Contract:   []byte(scAddress),
	}
	signature := signAdminTransaction(payload)
	body := SmartContractRequest{
		Payload:   payload,
		Signature: signature,
		IsView:    false,
	}
	requestBody, _ := json.Marshal(body)
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
	const SC_ENDPOINT = "api/smartContract"
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
		log.Warn("Could not parse int slice from this res - most likely empty response")
		return make([]int, 0), nil
	}
	return tokens, nil
}

func GetNftMetadata(user []byte, nftId int) (*string, error) {
	const SC_ENDPOINT = "api/smartContract"
	address := getServerAddress()
	url := fmt.Sprintf("%s/%s", address, SC_ENDPOINT)
	jsonAgs, _ := json.Marshal(nftId)
	scAddress := utils.GetConfigString(utils.NodeNftAddress)
	payload := SmartContractRequestPayload{
		Entrypoint: "_get_metadata",
		Args:       string(jsonAgs),
		Sender:     user,
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
		log.Warn("Metadata has crashed the node. Great.")
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
	return &responseBody.Output, nil
}

func GetNftMetadataParallel(user []byte, nfts []int) []string {
	var wg sync.WaitGroup
	metaChannel := make(chan string, len(nfts))
	for _, nft := range nfts {
		wg.Add(1)
		go func() {
			defer wg.Done()
			meta, err := GetNftMetadata(user, nft)
			if err == nil {
				metaChannel <- *meta
			} else {
				log.Warnf("Error while parallel fetching nft metadata: %s", err.Error())
			}
		}()
	}
	wg.Wait()
	close(metaChannel)
	var res []string
	for metadata := range metaChannel {
		res = append(res, metadata)
	}
	return res
}
