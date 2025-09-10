package transaction

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/utils"
	"encoding/hex"
	"time"
)

type TransactionDataDto struct {
	Fee       float64 `json:"fee" validate:"required,gt=0"`
	Amount    float64 `json:"amount" validate:"required,gt=0"`
	To        string  `json:"to" validate:"required,len=34"`
	Timestamp int64   `json:"timestamp" validate:"required,gt=0"`
	Message   string  `json:"message"`
}

type NewTransactionDto struct {
	Data   TransactionDataDto `json:"data" validate:"required"`
	Sig    string             `json:"sig" validate:"required,hexadecimal"`
	PubKey string             `json:"pubKey" validate:"required,hexadecimal"`
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

	pubKeyBytes, err := hex.DecodeString(tx.PubKey)
	if err != nil {
		return nil, apperror.BadRequest("Invalid public key format. It must be a valid hexadecimal string.", err)
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
			Message:   tx.Data.Message,
		},
		Sig:    sigBytes,
		PubKey: pubKeyBytes,
	}

	return parsed, nil
}
