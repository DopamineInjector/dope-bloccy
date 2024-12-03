package nft

import (
	"sync"
)

func GetUserNfts(userId []byte, ids []string, tokenIds []int) NftResponse {
	metadataEntries := getNftsMetadata(ids, tokenIds)
	nftEntries := getNftsImages(metadataEntries)
	return NftResponse{
		Nfts: nftEntries,
	}
}

func MintNft(request *MintNftRequest) (*NftMetadata, error) {
	metadata, err := mintNft(request.Description)
	if err != nil {
		return nil, err
	}
	err = postAvatar(metadata.ImageId, request.Image)
	if err != nil {
		return nil, err
	}
	return metadata, nil
}

func getNftsMetadata(ids []string, tokenIds []int) []NftMetaWithId {
	var wg sync.WaitGroup
	metaChannel := make(chan NftMetaWithId, len(ids))
	for i, id := range ids {
		wg.Add(1)
		go func() {
			defer wg.Done()
			meta, err := getNft(id)
			if err == nil {
				res := NftMetaWithId{
					Metadata: *meta,
					TokenId: tokenIds[i],
				}
				metaChannel <- res
			}
		}()
	}
	wg.Wait()
	close(metaChannel)
	var res []NftMetaWithId
	for metadata := range metaChannel {
		res = append(res, metadata)
	}
	return res
}

func getNftsImages(metadata []NftMetaWithId) []NftResponseEntry {
	var wg sync.WaitGroup
	outputChannel := make(chan NftResponseEntry, len(metadata))
	for _, meta := range metadata {
		wg.Add(1)
		go func() {
			defer wg.Done()
			image, err := getAvatar(meta.Metadata.ImageId)
			if err != nil {
				image = make([]byte, 0)
			}
			entry := NftResponseEntry{
				Id:          meta.Metadata.Id,
				Description: meta.Metadata.Description,
				Image:       image,
				TokenId: meta.TokenId,
			}
			outputChannel <- entry
		}()
	}
	wg.Wait()
	close(outputChannel)
	var res []NftResponseEntry
	for entry := range outputChannel {
		res = append(res, entry)
	}
	return res
}
