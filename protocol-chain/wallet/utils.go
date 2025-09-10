package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/mr-tron/base58"
	log "github.com/sirupsen/logrus"
)

func Base58Encode(input []byte) []byte {
	enCode := base58.Encode(input)

	return []byte(enCode)
}

func Base58Decode(input []byte) []byte {
	decode, err := base58.Decode(string(input))

	if err != nil {
		log.Panic(err)
	}

	return decode
}

func EncryptPrivateKeyForExport(privateKey string) (string, error) {
	padded := fmt.Sprintf("%s:::%s", conf.Wallet_Padding, privateKey)

	h := hmac.New(sha256.New, []byte(conf.SystemKey))
	_, err := h.Write([]byte(padded))

	if err != nil {
		return "", nil
	}

	hmacSum := h.Sum(nil)
	hmacHex := hex.EncodeToString(hmacSum)

	keyHash := sha256.Sum256([]byte(conf.SystemKey))
	key := keyHash[:]

	iv := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", nil
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", nil
	}

	ctWithTag := gcm.Seal(nil, iv, []byte(padded), nil)

	tag := ctWithTag[len(ctWithTag)-gcm.Overhead():]
	cipherOnly := ctWithTag[:len(ctWithTag)-gcm.Overhead()]

	out := append(append(iv, tag...), cipherOnly...)

	hexStr := hex.EncodeToString(out)

	return hexStr + ":::" + hmacHex, nil
}
