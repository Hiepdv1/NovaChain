package utils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"math/big"
)

func VerifyECDSASignature(pubKeyBytes, sigBytes []byte, data string) (bool, error) {
	hashData := sha256.Sum256([]byte(data))

	x := big.NewInt(0).SetBytes(pubKeyBytes[:32])
	y := big.NewInt(0).SetBytes(pubKeyBytes[32:])

	pub := &ecdsa.PublicKey{Curve: elliptic.P256(), X: x, Y: y}

	r := big.NewInt(0).SetBytes(sigBytes[:32])
	s := big.NewInt(0).SetBytes(sigBytes[32:])

	ok := ecdsa.Verify(pub, hashData[:], r, s)

	return ok, nil
}
