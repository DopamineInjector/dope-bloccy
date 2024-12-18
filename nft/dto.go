package nft

type NftMetadata struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	ImageId     string `json:"imageId"`
}

type NftMetaWithId struct {
	Metadata NftMetadata
	TokenId int `json:"tokenId"`
}

type PostMetadataDTO struct {
	Description string `json:"description"`
}

type NftResponseEntry struct {
	Id          string `json:"id"`
	TokenId			int `json:"tokenId"`
	Description string `json:"description"`
	Image       []byte `json:"image"`
}

type NftResponse struct {
	Nfts []NftResponseEntry `json:"nfts"`
}

type MintNftRequest struct {
	User        string `json:"user"`
	Description string `json:"description"`
	Image       []byte `json:"image"`
}

type TransferNftRequest struct {
	Sender string	`json:"sender"`
	Recipient string `json:"recipient"`
	TokenId int `json:"tokenId"`
}
