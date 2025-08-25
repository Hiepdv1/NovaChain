package gobtypes

import (
	"ChainServer/internal/common/ratelimiter"
	"encoding/gob"
)

func Init() {
	gob.Register(ratelimiter.BucketState{})
}
