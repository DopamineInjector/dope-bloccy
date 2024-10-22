package nft

import (
  log "github.com/sirupsen/logrus"
	"sync"
)

func GetNftsMetadata(ids []string) []NftMetadata {
  var wg sync.WaitGroup;
  metaChannel := make(chan NftMetadata, len(ids));
  for _, id := range ids {
    wg.Add(1);
    go func() {
      defer wg.Done();
      meta, err := getNft(id);
      log.Info(id);
      if err == nil {
        metaChannel <- *meta
      }
    }();
  }
  wg.Wait();
  close(metaChannel);
  var res []NftMetadata;
  for metadata := range metaChannel {
    res = append(res, metadata);
  }
  return res;
}
