package wallet

import (
	"ChainServer/internal/cache/redis"
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/env"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/utils"
	dbwallet "ChainServer/internal/db/wallet"
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
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

func (s *WalletService) CreateWallet(dto *dto.WalletParsed) (*string, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	errInternalCommon := apperror.Internal("Something went wrong. Please try again.", nil)

	existingWallet, err := s.dbRepo.ExistsWalletByPubKey(ctx, dto.PublicKey)

	if existingWallet {
		return nil, apperror.BadRequest("Wallet already exists", nil)
	}

	if err != nil {
		if apperr, ok := err.(*apperror.AppError); ok {
			return nil, apperr
		}
		log.Error("Failed to check if wallet exists: ", err)
		return nil, errInternalCommon
	}

	wallet, err := s.dbRepo.CreateWallet(ctx, dbwallet.CreateWalletParams{
		PublicKey:     hex.EncodeToString(dto.PublicKey),
		Address:       dto.Addr,
		PublicKeyHash: hex.EncodeToString(utils.PublicKeyHash(dto.PublicKey)),
		Balance:       fmt.Sprintf("%.8f", 0.000000000),
		CreateAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		LastLogin: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}, nil)

	if err != nil {
		log.Error("Failed to create wallet in database", err)
		return nil, errInternalCommon
	}

	payload := JWTWalletAuthPayload{
		ID:      wallet.ID.String(),
		Address: wallet.Address,
		Pubkey:  wallet.PublicKey,
	}

	token, err := utils.SignJWT(
		[]byte(env.Cfg.Jwt_Secret_Key),
		payload,
		time.Duration(env.Cfg.Jwt_TTL_Minutes)*time.Minute,
	)

	if err != nil {
		log.Error("JWT signing failed", err)
		return nil, errInternalCommon
	}

	err = redis.Set(
		ctx,
		helpers.BlacklistSigKey(helpers.AuthKeyTypeSig, hex.EncodeToString(dto.Sig)),
		token,
		time.Minute*time.Duration(env.Cfg.Wallet_Signature_Expiry_Minutes)+time.Minute,
	)

	if err != nil {
		log.Error("Redis set signature token failed", err)
		return nil, errInternalCommon
	}

	return &token, nil

}

func (s *WalletService) ImportWallet(dto *dto.WalletParsed) (*string, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	errInternalCommon := apperror.Internal("Something went wrong. Please try again.", nil)

	wallet, err := s.dbRepo.GetWalletByPubkey(ctx, dto.PublicKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			balance, err := s.rpcRepo.GetBalance(dto.Addr)

			if err != nil {
				log.Error("RPC balance lookup failed", err)
				return nil, errInternalCommon
			}

			if balance.Error != nil {
				return nil, apperror.BadRequest("Import wallet failed", errors.New(balance.Error.Message))
			}

			newWallet, err := s.dbRepo.CreateWallet(ctx, dbwallet.CreateWalletParams{
				Address:       dto.Addr,
				PublicKey:     string(dto.PublicKey),
				PublicKeyHash: string(utils.PublicKeyHash(dto.PublicKey)),
				Balance:       fmt.Sprintf("%.8f", balance.Balance),
				CreateAt: sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				},
				LastLogin: sql.NullTime{
					Time:  time.Now(),
					Valid: true,
				},
			}, nil)

			if err != nil {
				log.Error("Database insert wallet failed", err)
				return nil, errInternalCommon
			}

			wallet = &newWallet

		} else {
			log.Error("Database query wallet failed", err)
			return nil, errInternalCommon
		}
	}

	payload := JWTWalletAuthPayload{
		ID:      wallet.ID.String(),
		Address: wallet.Address,
		Pubkey:  wallet.PublicKey,
	}

	token, err := utils.SignJWT(
		[]byte(env.Cfg.Jwt_Secret_Key),
		payload,
		time.Duration(env.Cfg.Jwt_TTL_Minutes)*time.Minute,
	)

	if err != nil {
		log.Error("JWT signing failed", err)
		return nil, errInternalCommon
	}

	err = redis.Set(
		ctx,
		helpers.BlacklistSigKey(helpers.AuthKeyTypeSig, hex.EncodeToString(dto.Sig)),
		token,
		time.Minute*time.Duration(env.Cfg.Wallet_Signature_Expiry_Minutes)+time.Minute,
	)

	if err != nil {
		log.Error("Redis set signature token failed", err)
		return nil, errInternalCommon
	}

	return &token, nil
}

func (s *WalletService) GetWallet(pubkey []byte) (*dbwallet.Wallet, *apperror.AppError) {
	wallet, err := s.dbRepo.GetWalletByPubkey(context.Background(), pubkey)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound("Wallet not found", nil)
		}

		log.Error("Get wallet failed ", err)
		return nil, apperror.Internal("Something went wrong. Please try again.", nil)
	}

	return wallet, nil
}
