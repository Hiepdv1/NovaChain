package transaction

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/utils"
	"encoding/hex"
	"time"
)

type TransactionDataDto struct {
	Fee       float64 `json:"fee" validate:"required,gt=0"`
	Amount    float64 `json:"amount" validate:"required,gt=0"`
	To        string  `json:"to" validate:"required,len=34"`
	Timestamp int64   `json:"timestamp" validate:"required,gt=0"`
	Priority  uint64  `json:"priority" validate:"required,gte=0"`
}

type NewTransactionDto struct {
	Data TransactionDataDto `json:"data" validate:"required"`
	Sig  string             `json:"sig" validate:"required,hexadecimal"`
}

type SendTransactionDataDto struct {
	Amount          float64         `json:"amount" validate:"required,gt=0"`
	Fee             float64         `json:"fee" validate:"required,gt=0"`
	ReceiverAddress string          `json:"receiverAddress" validate:"required,len=34"`
	Priority        uint            `json:"priority" validate:"required,gte=0"`
	Transaction     dto.Transaction `json:"transaction" validate:"required"`
}

type SendTransactionDto struct {
	Data SendTransactionDataDto `json:"data" validate:"required"`
	Sig  string                 `json:"sig" validate:"required,hexadecimal"`
}

func (s *SendTransactionDto) ValidateAndParse() (any, error) {

	var inputs []TxInput
	var outputs []TxOutput

	txID, err := hex.DecodeString(s.Data.Transaction.ID)
	if err != nil {
		return nil, apperror.BadRequest("TxID is not format", nil)
	}

	if !utils.ValidateAddress(s.Data.ReceiverAddress) {
		return nil, apperror.BadRequest("Receiver address is invalid. It must be a valid blockchain address (34 characters, base58).", nil)
	}

	for _, in := range s.Data.Transaction.Inputs {
		inId, err := hex.DecodeString(in.ID)
		if err != nil {
			return nil, apperror.BadRequest("inID is not format", nil)
		}

		sigByte, err := hex.DecodeString(in.Signature)
		if err != nil {
			return nil, apperror.BadRequest("sig is not format", nil)
		}

		pubkey, err := hex.DecodeString(in.PubKey)
		if err != nil {
			return nil, apperror.BadRequest("pubkey is not format", nil)
		}

		input := TxInput{
			ID:        inId,
			Out:       in.Out,
			Signature: sigByte,
			PubKey:    pubkey,
		}
		inputs = append(inputs, input)
	}

	for _, out := range s.Data.Transaction.Outputs {
		pubKeyHash, err := hex.DecodeString(out.PubKeyHash)
		if err != nil {
			return nil, apperror.BadRequest("pubkey is not format", nil)
		}

		outputs = append(outputs, TxOutput{
			Value:      out.Value,
			PubKeyHash: pubKeyHash,
		})
	}

	tx := Transaction{
		ID:      txID,
		Inputs:  inputs,
		Outputs: outputs,
	}

	return SendTransactionDataParsed{
		Priority:     s.Data.Priority,
		Transaction:  tx,
		ReceiverAddr: s.Data.ReceiverAddress,
		Fee:          s.Data.Fee,
		Amount:       s.Data.Amount,
	}, nil

}

func (tx *NewTransactionDto) ValidateAndParse() (any, error) {
	now := time.Now().Unix()

	const maxDrift int64 = 180

	if tx.Data.Timestamp <= now {
		return nil, apperror.BadRequest("Timestamp is expired. It must be greater than the current Unix time (in seconds).", nil)
	}

	if tx.Data.Timestamp > now+maxDrift {
		return nil, apperror.BadRequest("Timestamp is too far in the future. It must not exceed 3 minutes ahead of the current Unix time (in seconds).", nil)
	}

	if !utils.ValidateAddress(tx.Data.To) {
		return nil, apperror.BadRequest("Recipient address is invalid. It must be a valid blockchain address (34 characters, base58).", nil)
	}

	sigBytes, err := hex.DecodeString(tx.Sig)
	if err != nil {
		return nil, apperror.BadRequest("Invalid signature format. It must be a valid hexadecimal string.", err)
	}

	toAddrBytes, err := utils.Base58Decode(tx.Data.To)

	if err != nil {
		return nil, apperror.BadRequest("Invalid format to address", nil)
	}

	parsed := &NewTransactionParsed{
		Data: TransactionDataParsed{
			Fee:       tx.Data.Fee,
			Amount:    tx.Data.Amount,
			To:        toAddrBytes,
			Timestamp: time.Unix(tx.Data.Timestamp, 0),
		},
		Sig: sigBytes,
	}

	return parsed, nil
}

type GetTransactionSearchDto struct {
	B_Hash          string `query:"b_hash" validate:"required,hexadecimal,len=64"`
	Search_Tx_Query string `query:"q" validate:"omitempty"`
	dto.PaginationQuery
}

type GetTransactionPendingDto struct {
	dto.PaginationQuery
}

type GetTransactionDetailDto struct {
	TxHash string `params:"tx_hash" validate:"required,hexadecimal,len=64"`
}
