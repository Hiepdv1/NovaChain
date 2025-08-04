package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"log"

	"golang.org/x/crypto/ripemd160"
)

var (
	version        = byte(0x00)
	checkSumlength = 4
)

func PublicKeyHash(pubKey []byte) []byte {
	pubHash := sha256.Sum256(pubKey)

	hasher := ripemd160.New()
	_, err := hasher.Write(pubHash[:])
	if err != nil {
		log.Panic(err)
	}

	publicRipMd := hasher.Sum(nil)
	return publicRipMd
}

func CheckSum(data []byte) []byte {
	firstHash := sha256.Sum256(data)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:checkSumlength]
}

func PubKeyToAddress(pubkey []byte) string {
	pubHash := PublicKeyHash(pubkey)
	versionedHash := append([]byte{version}, pubHash...)

	checksum := CheckSum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	return hex.EncodeToString(address)
}
