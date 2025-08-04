package wallet

import (
	"bytes"
	"core-blockchain/common/env"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"log"
	"math/big"

	"golang.org/x/crypto/ripemd160"
)

var conf = env.New()

var (
	checkSumlength = conf.WalletAddressCheckSum
	version        = byte(0x00)
)

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey  []byte
}

type WalletSerializable struct {
	PrivateKey []byte
	PublicKeyX []byte
	PublicKeyY []byte
	PublicKey  []byte
}

func NewWallet() *Wallet {
	private, public := NewKeyPair()
	return &Wallet{private, public}
}

func ValidateAddress(address string) bool {
	if len(address) != 34 {
		return false
	}

	fullHash := Base58Decode([]byte(address))

	checkSumFromHash := fullHash[int64(len(fullHash))-checkSumlength:]
	version := fullHash[0]
	pubKeyHash := fullHash[1 : int64(len(fullHash))-checkSumlength]
	checkSum := CheckSum(append([]byte{version}, pubKeyHash...))

	return bytes.Equal(checkSumFromHash, checkSum)
}

func (w *Wallet) Address() []byte {
	pubHash := PublicKeyHash(w.PublicKey)
	versionedHash := append([]byte{version}, pubHash...)

	checksum := CheckSum(versionedHash)

	fullHash := append(versionedHash, checksum...)
	address := Base58Encode(fullHash)

	return address
}

func (w *Wallet) Serialize() (*WalletSerializable, error) {

	privKey := w.PrivateKey.D.Bytes()
	pubKeyX := w.PrivateKey.PublicKey.X.Bytes()
	pubKeyY := w.PrivateKey.PublicKey.Y.Bytes()

	if pubKeyX == nil || pubKeyY == nil || privKey == nil {
		return nil, errors.New("missing private key")
	}

	return &WalletSerializable{
		PrivateKey: privKey,
		PublicKeyX: pubKeyX,
		PublicKeyY: pubKeyY,
		PublicKey:  w.PublicKey,
	}, nil
}

func (ws *WalletSerializable) Deserialize() (*Wallet, error) {
	priv := new(ecdsa.PrivateKey)
	priv.D = new(big.Int).SetBytes(ws.PrivateKey)
	priv.PublicKey.X = new(big.Int).SetBytes(ws.PublicKeyX)
	priv.PublicKey.Y = new(big.Int).SetBytes(ws.PublicKeyY)
	priv.PublicKey.Curve = elliptic.P256()

	return &Wallet{
		PrivateKey: *priv,
		PublicKey:  ws.PublicKey,
	}, nil
}

func NewKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()

	private, err := ecdsa.GenerateKey(curve, rand.Reader)

	if err != nil {
		log.Panic(err)
	}

	pub := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)

	return *private, pub

}

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
