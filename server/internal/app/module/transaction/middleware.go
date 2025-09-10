package transaction

import (
	"ChainServer/internal/common/response"
	"ChainServer/internal/common/utils"
	"encoding/hex"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func VerifyCreateTransactionSig(c *fiber.Ctx) error {
	var dto NewTransactionDto

	if err := c.BodyParser(&dto); err != nil {
		return response.Error(c, fiber.StatusBadRequest, "Invalid body", response.ErrBadRequest, err)
	}

	internalErr := response.Error(c, fiber.StatusInternalServerError, "Something went wrong. Please try again.", response.ErrInternal, nil)

	pubKeyBytes, err := hex.DecodeString(dto.PubKey)
	if err != nil {
		return internalErr
	}

	sigBytes, err := hex.DecodeString(dto.Sig)
	if err != nil {
		return internalErr
	}

	msgJson, err := json.Marshal(dto.Data)
	if err != nil {
		return internalErr
	}

	okSig, err := utils.VerifyECDSASignature(pubKeyBytes, sigBytes, string(msgJson))

	if err != nil {
		log.Errorf("VerifyCreateTransaction error: %v", err)
		return internalErr
	}

	if !okSig {
		return response.Error(
			c,
			fiber.StatusBadRequest,
			"Invalid signature",
			response.ErrBadRequest,
			nil,
		)
	}

	return c.Next()
}
