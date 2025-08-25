package middlewares

import (
	"ChainServer/internal/cache/redis"
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/ratelimiter"
	"ChainServer/internal/common/response"
	"crypto/md5"
	"fmt"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
	log "github.com/sirupsen/logrus"
)

type RateLimitConfig struct {
	GlobalLimiter ratelimiter.RateLimiter
	RouteLimiters map[string]ratelimiter.RateLimiter
	KeyFunc       func(c *fiber.Ctx) string
}

func RateLimitMiddleware(cfg RateLimitConfig) fiber.Handler {
	if cfg.GlobalLimiter == nil {
		log.Panic("Rate limiter instance is required")
	}

	if cfg.RouteLimiters == nil {
		cfg.RouteLimiters = make(map[string]ratelimiter.RateLimiter)
	}

	if cfg.KeyFunc == nil {
		cfg.KeyFunc = GetSecureKey
	}

	return func(c *fiber.Ctx) error {

		limiter := cfg.GlobalLimiter
		path := c.Path()
		for prefix, routeLimiter := range cfg.RouteLimiters {
			if strings.HasPrefix(path, prefix) {
				limiter = routeLimiter
				break
			}
		}

		key := redis.CacheKey{
			Namespace: redis.NamespaceRateLimit,
			Key:       cfg.KeyFunc(c),
		}

		_, result, err := limiter.Allow(c.Context(), key)
		if err != nil {
			apperr, ok := err.(*apperror.AppError)
			if !ok {
				apperr = apperror.Internal("Unexception error", err)
			}

			if apperr.ErrType == response.ErrTooManyRequests {
				c.Set("Retry-After", fmt.Sprintf("%.0f", result.RetryAfter.Seconds()))
				c.Set("X-Ratelimit-Remaining", fmt.Sprintf("%d", result.Remaining))
			}

			return response.Error(
				c,
				int(apperr.Status),
				apperr.Message,
				apperr.ErrType,
				nil,
			)
		}

		c.Set("X-Ratelimit-Remaining", fmt.Sprintf("%d", result.Remaining))
		return c.Next()
	}
}

func GetSecureKey(c *fiber.Ctx) string {
	userID, ok := c.Locals("user_id").(string)
	if ok && userID != "" {
		return userID
	}

	userAgent := c.Get("User-Agent")
	acceptLanguage := c.Get("Accept-Language")
	fingerprint := userAgent + "|" + acceptLanguage
	hash := md5.New()
	io.WriteString(hash, fingerprint)

	return fmt.Sprintf("fp:%x", hash.Sum(nil))
}
