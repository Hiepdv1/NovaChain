/* eslint-disable @typescript-eslint/no-explicit-any */
import { BaseErrorResponse, BaseResponse } from '@/shared/types/api';
import { useMutation } from '@tanstack/react-query';
import walletService from '../services/wallet.service';
import { WalletConnectPayload } from '../types/wallet';

export const useDisconnectWalletMutation = () => {
  return useMutation<BaseResponse<null>, BaseErrorResponse, null>({
    mutationFn: walletService.Disconnect,
  });
};

export const useWalletConnect = () => {
  return useMutation<
    BaseResponse<null>,
    BaseErrorResponse,
    WalletConnectPayload
  >({
    mutationFn: walletService.CreateWallet,
  });
};

export const useWalletImport = () => {
  return useMutation<
    BaseResponse<null>,
    BaseErrorResponse,
    WalletConnectPayload
  >({
    mutationFn: walletService.ImportWallet,
  });
};
