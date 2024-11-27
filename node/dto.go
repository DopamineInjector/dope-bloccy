package node;

type GetAccountInfoDto struct {
  WalletId []byte `json:"publicKey"`
}

type AccountInfoResponseDto struct {
  PublicKey string `json:"publicKey"`
  Balance float32 `json:"balance"`
}
