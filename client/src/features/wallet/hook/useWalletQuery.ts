/* eslint-disable @typescript-eslint/no-explicit-any */
import { BaseErrorResponse, BaseResponse } from '@/shared/types/api';
import { useMutation, useQuery } from '@tanstack/react-query';
import {
  Wallet,
  WalletConnectPayload,
  WalletQueryOptions,
} from '../types/wallet';
import walletService from '../services/wallet.service';

export const useWalletConnect = () => {
  return useMutation<
    BaseResponse<any>,
    BaseErrorResponse,
    WalletConnectPayload
  >({
    mutationFn: walletService.CreateWallet,
  });
};

export const useWalletImport = () => {
  return useMutation<
    BaseResponse<any>,
    BaseErrorResponse,
    WalletConnectPayload
  >({
    mutationFn: walletService.ImportWallet,
  });
};

export const useWalletQuery = (opts?: WalletQueryOptions) => {
  return useQuery<Wallet, BaseErrorResponse>({
    queryKey: ['wallet'],
    queryFn: async () => {
      const res = await walletService.GetWallet();
      return res.data;
    },
    ...opts,
  });
};

export const useDisconnectWalletMutation = () => {
  return useMutation<BaseResponse<null>, BaseErrorResponse, null>({
    mutationFn: walletService.Disconnect,
  });
};
