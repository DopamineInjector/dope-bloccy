package nft

import (
	"sync"
)

func GetUserNfts(userId []byte) NftResponse {
  // TODO - actually request a node to give info about user nfts
  ids := [2]string{"6a389c77-6bab-4e0c-98b6-4c25ca23c60e", "2bc027ad-573a-4a58-a049-9d6c04800e5c"};
  metadataEntries := getNftsMetadata(ids[:]);
  nftEntries := getNftsImages(metadataEntries);
  return NftResponse{
    Nfts: nftEntries,
  };
}

func MintNft(request *MintNftRequest) error {
  metadata, err := mintNft(request.Description);
  if err != nil {
    return err
  }
  err = postAvatar(metadata.ImageId, request.Image);
  if err != nil {
    return err
  }
  // TODO - actually ask node to mint the nft
  return nil
}

func getNftsMetadata(ids []string) []NftMetadata {
  var wg sync.WaitGroup;
  metaChannel := make(chan NftMetadata, len(ids));
  for _, id := range ids {
    wg.Add(1);
    go func() {
      defer wg.Done();
      meta, err := getNft(id);
      if err == nil {
        metaChannel <- *meta
      }
    }();
  }
  wg.Wait();
  close(metaChannel);
  var res []NftMetadata;
  for metadata := range metaChannel {
    res = append(res, metadata)
  }
  return res
}

func getNftsImages(metadata []NftMetadata) []NftResponseEntry {
  var wg sync.WaitGroup;
  outputChannel := make(chan NftResponseEntry, len(metadata));
  for _, meta := range metadata {
    wg.Add(1);
    go func() {
      defer wg.Done();
      image, err := getAvatar(meta.ImageId);
      if err != nil {
        image = make([]byte, 0);
      }
      entry := NftResponseEntry{
        Id: meta.Id,
        Description: meta.Description,
        Image: image,
      }
      outputChannel <- entry
    }()
  }
  wg.Wait();
  close(outputChannel);
  var res []NftResponseEntry;
  for entry := range outputChannel {
    res = append(res, entry)
  }
  return res
}
