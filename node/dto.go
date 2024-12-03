package node

type GetAccountInfoDto struct {
	WalletId []byte `json:"publicKey"`
}

type AccountInfoResponseDto struct {
	PublicKey string  `json:"publicKey"`
	Balance   float32 `json:"balance"`
}

type TransferRequestPayload struct {
	Sender    []byte  `json:"sender"`
	Recipient []byte  `json:"recipient"`
	Amount    float32 `json:"amount"`
}

type TransferRequest struct {
	Payload   TransferRequestPayload `json:"payload"`
	Signature []byte                 `json:"signature"`
}

type SmartContractRequestPayload struct {
	Sender     []byte `json:"sender"`
	Contract   []byte `json:"contract"`
	Entrypoint string `json:"entrypoint"`
	Args       string `json:"args"`
}

type SmartContractRequest struct {
	Payload   SmartContractRequestPayload `json:"payload"`
	Signature []byte                      `json:"signature"`
	IsView    bool                        `json:"view"`
}

type SmartContractResponse struct {
	Output string `json:"output"`
}

type MintNftArgs struct {
	Recipient   []byte `json:"owner"`
	MetadataUri string `json:"metadata_uri"`
}

type OwnedByArgs struct {
	Owner []byte `json:"owner"`
}
