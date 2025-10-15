package transaction

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/utils"
	"encoding/hex"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func VerifyCreateSignatureMiddleware(c *fiber.Ctx) error {
	dto, apperr := helpers.GetLocalRaw[NewTransactionDto](c, "body")
	if apperr != nil {
		return apperr.Response(c)
	}

	internalErr := apperror.Internal("Unable to verify transaction signature. Please try again later.", nil)

	walletPayload, apperr := helpers.GetLocalWallet(c)
	if apperr != nil {
		return apperr
	}

	pubKeyBytes, err := hex.DecodeString(walletPayload.Data.Pubkey)
	if err != nil {
		log.Errorf("Decode public key error: %v", err)
		return internalErr.Response(c)
	}

	sigBytes, err := hex.DecodeString(dto.Sig)
	if err != nil {
		log.Errorf("Decode signature error: %v", err)
		return internalErr.Response(c)
	}

	msgJson, err := json.Marshal(dto.Data)
	if err != nil {
		log.Errorf("Marshal transaction data error: %v", err)
		return internalErr.Response(c)
	}

	okSig, err := utils.VerifyECDSASignature(pubKeyBytes, sigBytes, string(msgJson))
	if err != nil {
		log.Errorf("Signature verification process error: %v", err)
		return internalErr.Response(c)
	}

	if !okSig {
		return apperror.BadRequest("The provided transaction signature is invalid.", nil).Response(c)
	}

	return c.Next()
}

func VerifySendPayloadMiddleware(c *fiber.Ctx) error {
	wallet, apperr := helpers.GetLocalWallet(c)
	if apperr != nil {
		return apperr.Response(c)
	}

	dto, apperr := helpers.GetLocalRaw[SendTransactionDto](c, "body")
	if apperr != nil {
		return apperr.Response(c)
	}

	pubKeyBytes, err := hex.DecodeString(wallet.Data.Pubkey)
	if err != nil {
		return apperror.Internal("Failed to decode public key.", err).Response(c)
	}

	sigBytes, err := hex.DecodeString(dto.Sig)
	if err != nil {
		return apperror.Internal("Failed to decode payload signature.", err).Response(c)
	}

	dataToSign, err := json.Marshal(dto.Data)
	if err != nil {
		return apperror.Internal("Failed to serialize payload data.", err).Response(c)
	}

	okSig, err := utils.VerifyECDSASignature(pubKeyBytes, sigBytes, string(dataToSign))
	if err != nil {
		log.Errorf("Payload signature verification error: %v", err)
		return apperror.Internal("Error occurred during payload signature verification.", err).Response(c)
	}

	if !okSig {
		return apperror.BadRequest("Invalid payload signature.", nil).Response(c)
	}

	return c.Next()
}
