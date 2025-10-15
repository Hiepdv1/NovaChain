package middlewares

import (
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/env"
	"ChainServer/internal/common/utils"
	"encoding/base64"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

func DecryptBodyMiddleware(field *string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var encryptedStr string

		if field != nil {
			var bodyMap map[string]any
			if err := json.Unmarshal(c.Body(), &bodyMap); err != nil {
				return fiber.NewError(fiber.StatusBadRequest, "invalid body")
			}
			enc, ok := bodyMap[*field].(string)
			if !ok || enc == "" {
				return apperror.BadRequest("payload missing", nil).Response(c)
			}
			encryptedStr = enc
		} else {
			encryptedStr = string(c.Body())
		}

		data, err := base64.StdEncoding.DecodeString(encryptedStr)
		if err != nil {
			return apperror.BadRequest("invalid base64 string", nil).Response(c)
		}

		dataBytes, err := utils.DecryptData(data, env.Cfg.Encode_data_secret_Key)
		if err != nil {
			log.Error("failed to decrypt data: ", err)
			return apperror.Internal("failed to decrypt data", nil).Response(c)
		}

		c.Request().Header.SetContentType("application/json")
		c.Request().SetBody(dataBytes)

		return c.Next()
	}
}
