package dto

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/env"
	"ChainServer/internal/common/utils"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type WalletAuthData struct {
	Nonce     string `json:"nonce" validate:"required,uuid4"`
	PublicKey string `json:"publickey" validate:"required,hexadecimal,len=128"`
	Timestamp int64  `json:"timestamp" validate:"required,gt=0"`
	Address   string `json:"address" validate:"required"`
}

type WalletRequest struct {
	Data WalletAuthData `json:"data" validate:"required"`
	Sig  string         `json:"sig" validate:"required"`
}

type WalletParsed struct {
	Timestamp int64
	PublicKey []byte
	Nonce     uuid.UUID
	Sig       []byte
	Addr      string
}

func (w *WalletRequest) ValidateAndParse() (any, error) {

	now := time.Now().Unix()
	if w.Data.Timestamp <= 0 {
		return nil, apperror.BadRequest("Missing or invalid timestamp value", nil)
	}

	maxSkewSeconds := env.Cfg.Wallet_Signature_Expiry_Minutes * int64(time.Minute)

	if w.Data.Timestamp < now-maxSkewSeconds {
		return nil, apperror.BadRequest(
			"Timestamp expired: must be within the last %d minutes (server time in seconds: %d, expected seconds not milliseconds)",
			nil,
		)
	}

	if w.Data.Timestamp > now+maxSkewSeconds {
		return nil, apperror.BadRequest(
			fmt.Sprintf("Timestamp is too far in the future: must be within %d minutes ahead of server time ( server time by seconds: %d )", maxSkewSeconds/60, now),
			nil,
		)
	}

	nonceUUID, err := uuid.Parse(w.Data.Nonce)
	if err != nil {
		return nil, apperror.BadRequest("Invalid nonce format, must be UUID", err)
	}

	pubBytes, err := hex.DecodeString(w.Data.PublicKey)
	if err != nil {
		return nil, apperror.BadRequest("Invalid public key format, must be hex", err)
	}

	addrBytes := utils.PubKeyToAddress(pubBytes)

	if string(addrBytes) != w.Data.Address {
		return nil, apperror.BadRequest(
			"Address mismatch: provided address does not match the public key"+fmt.Sprintf("expected %s, got %s", addrBytes, w.Data.Address),
			nil,
		)
	}

	if !utils.ValidateAddress(w.Data.Address) {
		return nil, apperror.BadRequest("Invalid address format", nil)
	}

	sigBytes, err := hex.DecodeString(w.Sig)
	if err != nil {
		return nil, apperror.BadRequest("Invalid signature format, must be hex", err)
	}

	return &WalletParsed{
		Timestamp: w.Data.Timestamp,
		PublicKey: pubBytes,
		Nonce:     nonceUUID,
		Sig:       sigBytes,
		Addr:      w.Data.Address,
	}, nil
}
