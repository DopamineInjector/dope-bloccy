package nft

type NftMetadata struct {
	Id          string `json:"id"`
	Description string `json:"description"`
	ImageId     string `json:"imageId"`
}

type PostMetadataDTO struct {
	Description string `json:"description"`
}

type NftResponseEntry struct {
	Id          string `json:"id"`
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
