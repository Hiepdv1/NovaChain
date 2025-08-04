package wallet

import (
	"bytes"
	"core-blockchain/common/utils"
	"crypto/elliptic"
	"encoding/gob"
	"errors"
	"fmt"

	"os"
	"path"
	"path/filepath"
	"runtime"
)

var (
	_, file, _, _ = runtime.Caller(0)

	Root           = filepath.Join(filepath.Dir(file), "../")
	walletPath     = path.Join(Root, "/.chain/")
	walletFileName = ".wallets"
)

type WalletPool struct {
	Wallets map[string]*Wallet
}

type WalletPoolSerializable struct {
	Wallets map[string]*WalletSerializable
}

func ChainExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func InitializeWallets(cwd bool) (*WalletPool, error) {
	var walletPool WalletPool

	err := walletPool.LoadFile(cwd)

	return &walletPool, err
}

func (wp *WalletPool) GetWallet(address string) (Wallet, error) {
	var wallet *Wallet
	var ok bool
	w := *wp
	if wallet, ok = w.Wallets[address]; !ok {
		return *new(Wallet), errors.New("invalid address")
	}

	return *wallet, nil
}

func (wp *WalletPool) AddWallet() string {
	wallet := NewWallet()
	address := string(wallet.Address())

	wp.Wallets[address] = wallet

	return address
}

func (wp *WalletPool) GetAllAddress() []string {
	var addresses []string
	for address := range wp.Wallets {
		addresses = append(addresses, address)
	}

	return addresses
}

func (wp *WalletPool) LoadFile(cwd bool) error {
	if !ChainExists(walletPath) {
		err := os.MkdirAll(walletPath, 0755)
		if err != nil {
			return err
		}
		fmt.Println(".chain directory created:", walletPath)
	}

	walletFile := path.Join(walletPath, walletFileName)

	if cwd {
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		walletFile = path.Join(dir, walletFileName)
	}

	if _, err := os.Stat(walletFile); os.IsNotExist(err) {
		file, err := os.Create(walletFile)
		if err != nil {
			return err
		}
		file.Close()
		fmt.Println(".wallets file created:", walletPath)
	}

	walletPool := WalletPoolSerializable{
		Wallets: map[string]*WalletSerializable{},
	}
	fileContent, err := os.ReadFile(walletFile)
	if err != nil {
		return err
	}

	if len(fileContent) == 0 {
		wp.Wallets = make(map[string]*Wallet)
		return nil
	}

	gob.Register(elliptic.P256())
	decoder := gob.NewDecoder(bytes.NewReader(fileContent))
	err = decoder.Decode(&walletPool)

	if err != nil {
		return err
	}

	ws := WalletPool{
		Wallets: map[string]*Wallet{},
	}

	for addr, wallet := range walletPool.Wallets {
		w, err := wallet.Deserialize()
		if err != nil {
			return err
		}
		ws.Wallets[addr] = w
	}

	wp.Wallets = ws.Wallets

	return nil
}

func (wp *WalletPool) SaveFile(cwd bool) {
	walletFile := path.Join(walletPath, walletFileName)

	if cwd {
		dir, err := os.Getwd()
		utils.ErrorHandle(err)
		walletFile = path.Join(dir, walletFileName)
	}

	var content bytes.Buffer
	gob.Register(elliptic.P256())

	walletPool := WalletPoolSerializable{
		Wallets: map[string]*WalletSerializable{},
	}

	for addr, wallet := range wp.Wallets {
		w, err := wallet.Serialize()
		utils.ErrorHandle(err)
		walletPool.Wallets[addr] = w
	}

	encoder := gob.NewEncoder(&content)
	err := encoder.Encode(&walletPool)
	utils.ErrorHandle(err)

	err = os.WriteFile(walletFile, content.Bytes(), 0644)
	utils.ErrorHandle(err)
}
