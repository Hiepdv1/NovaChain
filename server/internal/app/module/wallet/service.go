package wallet

import (
	redisdb "ChainServer/internal/cache/redis"
	"ChainServer/internal/common/apperror"
	"ChainServer/internal/common/dto"
	"ChainServer/internal/common/env"
	"ChainServer/internal/common/helpers"
	"ChainServer/internal/common/types"
	"ChainServer/internal/common/utils"
	dbwallet "ChainServer/internal/db/wallet"
	"context"
	"database/sql"
	"encoding/hex"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

type WalletService struct {
	rpcRepo   RPCWalletRepository
	dbRepo    DBWalletRepository
	cacheRepo CacheWalletRepository
}

func NewWalletService(rpcRepo RPCWalletRepository, dbRepo DBWalletRepository, cacheRepo CacheWalletRepository) *WalletService {
	return &WalletService{
		rpcRepo:   rpcRepo,
		dbRepo:    dbRepo,
		cacheRepo: cacheRepo,
	}
}

func (s *WalletService) CreateWallet(dto *dto.WalletParsed) (*string, *apperror.AppError) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	errInternalCommon := apperror.Internal("Something went wrong. Please try again.", nil)

	existingWallet, err := s.dbRepo.ExistsWalletByPubKey(ctx, dto.PublicKey, nil)

	if existingWallet {
		return nil, apperror.BadRequest("Wallet already exists", nil)
	}

	if err != nil {
		log.Error("Failed to check if wallet exists: ", err)
		return nil, errInternalCommon
	}

	wallet, err := s.dbRepo.CreateWallet(ctx, dbwallet.CreateWalletParams{
		PublicKey:     helpers.StringToNullString(hex.EncodeToString(dto.PublicKey)),
		Address:       helpers.StringToNullString(dto.Addr),
		PublicKeyHash: hex.EncodeToString(utils.PublicKeyHash(dto.PublicKey)),
		Balance:       "0",
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

	payload := types.JWTWalletAuthPayload{
		ID:      wallet.ID.String(),
		Address: wallet.Address.String,
		Pubkey:  wallet.PublicKey.String,
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

	err = redisdb.Set(
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

	pubKeyHash := hex.EncodeToString(utils.PublicKeyHash(dto.PublicKey))

	wallet, err := s.dbRepo.GetWalletByPubKeyHash(ctx, pubKeyHash, nil)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			balance, err := s.rpcRepo.GetBalance(dto.Addr)

			if err != nil {
				log.Error("RPC balance lookup failed", err)
				return nil, errInternalCommon
			}

			if balance.Error != nil {
				log.Error(balance.Error.Message)
				return nil, apperror.BadRequest("Import wallet failed", errors.New(balance.Error.Message))
			}

			newWallet, err := s.dbRepo.CreateWallet(ctx, dbwallet.CreateWalletParams{
				PublicKey:     helpers.StringToNullString(hex.EncodeToString(dto.PublicKey)),
				Address:       helpers.StringToNullString(dto.Addr),
				PublicKeyHash: pubKeyHash,
				Balance:       helpers.FormatDecimal(balance.Balance, 8),
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

	if !wallet.Address.Valid || !wallet.PublicKey.Valid {
		walletUpdated, err := s.dbRepo.UpdateWalletInfoByWalletID(ctx, dbwallet.UpdateWalletInfoByWalletIDParams{
			PublicKey:     helpers.StringToNullString(hex.EncodeToString(dto.PublicKey)),
			Address:       helpers.StringToNullString(dto.Addr),
			PublicKeyHash: helpers.StringToNullString(pubKeyHash),
			Balance:       helpers.StringToNullString(""),
			ID:            wallet.ID,
		}, nil)

		if err != nil {
			log.Error(err)
			return nil, errInternalCommon
		}

		wallet = walletUpdated
	}

	payload := types.JWTWalletAuthPayload{
		ID:      wallet.ID.String(),
		Address: wallet.Address.String,
		Pubkey:  wallet.PublicKey.String,
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

	err = redisdb.Set(
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

func (s *WalletService) Disconnect(token string, payload utils.JWTPayload[types.JWTWalletAuthPayload]) {
	ttl := time.Until(payload.ExpiresAt.Time) + 2*time.Minute

	if ttl <= 0 {
		ttl = 2 * time.Minute
	}

	err := redisdb.Set(context.Background(),
		helpers.BlacklistSigKey(helpers.AuthKeyTypeJWT, token),
		"blacklisted",
		ttl,
	)
	if err != nil {
		log.Errorf("Failed set blacklist token %s: %v", token, err)
	}
}

func (s *WalletService) GetWallet(pubkey []byte) (*dbwallet.Wallet, *apperror.AppError) {
	ctx := context.Background()

	key := hex.EncodeToString(pubkey)

	wallet, err := s.dbRepo.GetWalletByPubkey(ctx, pubkey, nil)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.NotFound("Wallet not found", nil)
		}

		log.Error("Get wallet failed ", err)
		return nil, apperror.Internal("Something went wrong. Please try again.", nil)
	}

	if err != nil {
		log.Errorf("Failed to set wallet cache with pubKey %s: %v", key, err)
	}

	return wallet, nil
}
