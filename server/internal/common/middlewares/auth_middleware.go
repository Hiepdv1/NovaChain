package middlewares

import (
	"ChainServer/internal/cache/redis"
	"ChainServer/internal/common/constants"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/env"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/response"
	"ChainServer/internal/common/utils"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func VerifyWalletSignature(c *fiber.Ctx) error {
	val := c.Locals("body")
	if val == nil {
		return response.Error(c, fiber.StatusBadRequest, "missing parsed body", response.ErrBadRequest, nil)
	}

	parsed, ok := val.(*dto.WalletParsed)
	if !ok {
		return response.Error(c, fiber.StatusBadRequest, "invalid parsed body type", response.ErrBadRequest, nil)
	}

	sigHex := hex.EncodeToString(parsed.Sig)

	exists, err := redis.Exists(context.Background(), helpers.BlacklistSigKey(helpers.AuthKeyTypeSig, sigHex))

	if err != nil {
		log.Error("Failed to check signature existence in Redis", err)
		return response.Error(c, fiber.StatusInternalServerError, "Signature verification failed. Please try again.", response.ErrInternal, nil)
	}

	if exists {
		return response.Error(c, fiber.StatusUnauthorized, "Signature has been revoked", response.ErrUnauthorized, nil)
	}

	expiry := time.UnixMilli(parsed.Timestamp * 1000).Add(time.Duration(env.Cfg.Wallet_Signature_Expiry_Minutes) * time.Minute)
	if time.Now().After(expiry) {
		return response.Error(
			c,
			fiber.StatusUnauthorized,
			"Signature expired",
			response.ErrUnauthorized,
			nil,
		)
	}

	msg := dto.WalletAuthData{
		Nonce:     parsed.Nonce.String(),
		PublicKey: hex.EncodeToString(parsed.PublicKey),
		Timestamp: parsed.Timestamp,
		Address:   parsed.Addr,
	}

	msgJson, err := json.Marshal(msg)

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "failed to marshal message", response.ErrInternal, err)
	}

	okSig, err := utils.VerifyECDSASignature(parsed.PublicKey, parsed.Sig, string(msgJson))

	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, "Internal error", response.ErrInternal, err)
	}

	if !okSig {
		return response.Error(c, fiber.StatusUnauthorized, "Invalid signature", response.ErrUnauthorized, nil)
	}

	return c.Next()
}

func JWTAuthMiddleware[T any](c *fiber.Ctx) error {
	token := c.Cookies(constants.CookieAccessToken)
	if token == "" {
		return response.Error(
			c,
			fiber.StatusUnauthorized,
			"Access denied: missing token",
			response.ErrUnauthorized,
			nil,
		)
	}

	exists, err := redis.Exists(
		context.Background(),
		helpers.BlacklistSigKey(helpers.AuthKeyTypeJWT, token),
	)

	if err != nil {
		log.Errorf("Check blacklist token error: %v", err)

		return response.Error(
			c,
			fiber.StatusInternalServerError,
			"Server error: failed to validate token",
			response.ErrInternal,
			nil,
		)
	}

	if exists {
		return response.Error(
			c,
			fiber.StatusUnauthorized,
			"Access denied: token has been revoked",
			response.ErrUnauthorized,
			nil,
		)
	}

	payload, err := utils.VerifyJWT[T](
		[]byte(env.Cfg.Jwt_Secret_Key),
		token,
	)

	if err != nil {
		return response.Error(
			c,
			fiber.StatusUnauthorized,
			fmt.Sprintf("Invalid or expired token: %s", err.Error()),
			response.ErrUnauthorized,
			err,
		)
	}

	helpers.SetLocalWallet(c, payload)

	return c.Next()
}
