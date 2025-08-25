package response

import "ChainServer/internal/common/env"

func GetMessage(code ErrorType, rawMsg string) string {
	if env.Cfg.AppEnv == "production" {
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
