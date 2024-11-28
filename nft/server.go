package nft

import (
	"bytes"
	"dope-bloccy/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const METADATA_ENDPOINT = "metadata"
const AVATAR_ENDPOINT = "avatars"

func getServerAddress() string {
	return utils.GetConfigString(utils.MetadataServer)
}

func getNft(id string) (*NftMetadata, error) {
	address := getServerAddress()
	url := fmt.Sprintf("%s/%s/%s", address, METADATA_ENDPOINT, id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	res := &NftMetadata{}
	err = json.Unmarshal(body, res)
	return res, err
}

func mintNft(description string) (*NftMetadata, error) {
	address := getServerAddress()
	url := fmt.Sprintf("%s/%s", address, METADATA_ENDPOINT)
	data := PostMetadataDTO{
		Description: description,
	}
	requestBody, _ := json.Marshal(data)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("Error while minting metadata on nft server")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	metadata := &NftMetadata{}
	err = json.Unmarshal(body, metadata)
	return metadata, err
}

func getAvatar(id string) ([]byte, error) {
	address := getServerAddress()
	url := fmt.Sprintf("%s/%s/%s", address, AVATAR_ENDPOINT, id)
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	return body, err
}

func postAvatar(id string, avatar []byte) error {
	address := getServerAddress()
	url := fmt.Sprintf("%s/%s/%s", address, AVATAR_ENDPOINT, id)
	resp, err := http.Post(url, "image/png", bytes.NewBuffer(avatar))
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Error while creating avatar in metadata server")
	}
	return nil
}
