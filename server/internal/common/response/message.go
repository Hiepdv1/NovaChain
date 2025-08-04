package response

import (
	"ChainServer/internal/common/env"
)

var config = env.New()

func GetMessage(code ErrorCode, rawMsg string) string {
	if config.AppEnv == "production" {
		switch code {
		case ErrInternal:
			return "Something went wrong, please try again later"
		case ErrBadRequest:
			return "Invalid request"
		case ErrNotFound:
			return "Not Found"
		default:
			return "An error occurred"
		}
	}

	return rawMsg
}
