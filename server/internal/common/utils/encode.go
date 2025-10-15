package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/mr-tron/base58"
)

func Base58Encode(input []byte) []byte {
	enCode := base58.Encode(input)

	return []byte(enCode)
}

func Base58Decode(input string) ([]byte, error) {
	decode, err := base58.Decode(input)

	if err != nil {
		return nil, err
	}

	return decode, nil
}

func DoubleSHA256(data []byte) []byte {
	first := sha256.Sum256(data)
	second := sha256.Sum256(first[:])
	return second[:]
}

func normalizeKey(secret []byte) []byte {
	hash := sha256.Sum256(secret)
	return hash[:]
}

func EncryptData(data, secretKey []byte) (string, error) {
	key := normalizeKey(secretKey)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, nonce, data, nil)

	tagSize := 16
	if len(ciphertext) < tagSize {
		return "", fmt.Errorf("ciphertext too short for tag")
	}
	tag := ciphertext[len(ciphertext)-tagSize:]
	encData := ciphertext[:len(ciphertext)-tagSize]

	buffer := append(append(nonce, tag...), encData...)

	return base64.StdEncoding.EncodeToString(buffer), nil
}

func DecryptData(encryptedBytes, secretKey []byte) ([]byte, error) {
	key := normalizeKey(secretKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := 12
	tagSize := 16

	if len(encryptedBytes) < nonceSize+tagSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := encryptedBytes[:nonceSize]
	tag := encryptedBytes[nonceSize : nonceSize+tagSize]
	ciphertext := encryptedBytes[nonceSize+tagSize:]

	ciphertextWithTag := append(ciphertext, tag...)

	plaintext, err := aesGCM.Open(nil, nonce, ciphertextWithTag, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
