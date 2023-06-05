package context

import "github.com/ONBUFF-IP-TOKEN/baseapp/base"

type ReqPostWalletETHResult struct {
	PlatformID     int64  `json:"platform_id"`
	BlockNumber    string `json:"block_number"`
	Protocol       string `json:"protocol"`
	Contract       string `json:"contract"`
	Symbol         string `json:"symbol"`
	TxHash         string `json:"tx_hash"`
	Method         string `json:"method"`
	FromAddr       string `json:"from_addr"`
	ToAddr         string `json:"to_addr"`
	TokenID        string `json:"token_id"` // nft 전송에만 유효
	Value          string `json:"value"`
	Nonce          uint64 `json:"nonce"`
	TransactionFee uint64 `json:"transaction_fee"`
	GasPrice       uint64 `json:"gas_price"`
	GasUsed        uint64 `json:"gas_used"`
	Status         string `json:"status"` // success or fail
	CreateAt       uint64 `json:"create_at"`
}

func NewReqPostWalletETHResult() *ReqPostWalletETHResult {
	return new(ReqPostWalletETHResult)
}

func (o *ReqPostWalletETHResult) CheckValidate() *base.BaseResponse {

	return nil
}

// sui 계열 정보
type SUI_CB_GasUsed struct {
	Owner                   string `json:"owenr"`
	ComputationCost         string `json:"computation_cost"`
	StorageCost             string `json:"storage_cost"`
	StorageRebate           string `json:"storage_rebate"`
	NonRefundableStorageFee string `json:"non_refundable_storagefee"`
}

type SUI_CB_Balance_Change struct {
	Owner    string `json:"owner"`
	ObjectID string `json:"object_id"`
	Symbol   string `json:"symbol"`
	Amount   string `json:"amount"`
}

type SUI_CB_Balance_Changes struct {
	Status     string                   `json:"status"`
	PlatformID int64                    `json:"platform_id"`
	CheckPoint string                   `json:"check_point"`
	Digest     string                   `json:"digest"`
	Sender     string                   `json:"sender"`
	Balances   []*SUI_CB_Balance_Change `json:"balances"`
	GasUsed    SUI_CB_GasUsed           `json:"gas_used"`
}

func NewSUI_CB_Balance_Changes() *SUI_CB_Balance_Changes {
	return new(SUI_CB_Balance_Changes)
}

func (o *SUI_CB_Balance_Changes) CheckValidate() *base.BaseResponse {
	return nil
}

// nft 콜백 공통 정보

// nft mint 콜백
type SUI_CB_NFT_Minted_Json struct {
	Creator  string `json:"creator" mapstructure:"creator"`
	Name     string `json:"name" mapstructure:"name"`
	ObjectID string `json:"object_id" mapstructure:"object_id"`
}

// nft 전송 콜백
type SUI_CB_NFT_Transfer_Json struct {
	From     string `json:"from" mapstructure:"from"`
	Name     string `json:"name" mapstructure:"name"`
	ObjectID string `json:"object_id" mapstructure:"object_id"`
	To       string `json:"to" mapstructure:"to"`
}

// nft burn 콜백
type SUI_CB_NFT_Burn_Json struct {
	Name     string `json:"name" mapstructure:"name"`
	ObjectID string `json:"object_id" mapstructure:"object_id"`
	Owner    string `json:"owner" mapstructure:"owner"`
}
type SUI_CB_NFT struct {
	Status         string                  `json:"status"`
	ErrorMessage   string                  `json:"error_message"`
	Sender         string                  `json:"sender"`
	Digest         string                  `json:"digest"`
	NFTObjectID    string                  `json:"nft_obejct_id"`
	Module         string                  `json:"module"`
	Function       string                  `json:"function"`
	BCS            string                  `json:"bcs"`
	BalanceChanges []SUI_CB_Balance_Change `json:"balance_changes"`
	GasUsed        SUI_CB_GasUsed          `json:"gas_used"`

	MintInfo    SUI_CB_NFT_Minted_Json   `json:"nft_mint_info"`
	TranferInfo SUI_CB_NFT_Transfer_Json `json:"nft_transfer_info"`
	BurnInfo    SUI_CB_NFT_Burn_Json     `json:"nft_burn_info"`
}

func NewSUI_CB_NFT() *SUI_CB_NFT {
	return new(SUI_CB_NFT)
}

func (o *SUI_CB_NFT) CheckValidate() *base.BaseResponse {
	return nil
}
