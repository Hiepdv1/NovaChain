package wallet

import (
	"ChainServer/internal/common/client"
	"encoding/json"
	"errors"
)

type WalletService struct {
	rpcRepo RPCWalletRepository
	dbRepo  DBWalletRepository
}

func NewWalletService(rpcRepo RPCWalletRepository, dbRepo DBWalletRepository) *WalletService {
	return &WalletService{
		rpcRepo: rpcRepo,
		dbRepo:  dbRepo,
	}
}

func (s *WalletService) GetBalance(address string) (*Balance, error) {
	data, err := s.rpcRepo.GetBalance(address)

	if err != nil {
		return nil, err
	}

	var rpcResp client.RPCResponse
	if err := json.Unmarshal(data, &rpcResp); err != nil {
		panic(err)
	}

	var balance Balance
	if err := json.Unmarshal(rpcResp.Result, &balance); err != nil {
		panic(err)
	}

	if balance.Error != nil && balance.Error.Code != 0 {
		return nil, errors.New("failed to: " + balance.Error.Message)
	}

	balance.Error = nil

	return &balance, nil
}
