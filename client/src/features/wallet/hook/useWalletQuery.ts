import {
  BaseErrorResponse,
  BaseResponse,
  BaseResponseList,
} from '@/shared/types/api';
import { useQuery } from '@tanstack/react-query';
import { RecentTransaction, TxSumary, Wallet } from '../types/wallet';
import walletService from '../services/wallet.service';
import { PaginationParam, QueryOptions } from '@/shared/types/query';

export const useWalletQuery = (opts?: QueryOptions<BaseResponse<Wallet>>) => {
  return useQuery<BaseResponse<Wallet>, BaseErrorResponse>({
    queryKey: ['wallet'],
    queryFn: async () => {
      const res = await walletService.GetWallet();
      return res;
    },
    ...opts,
  });
};

export const useWalletSummary = (
  otps?: QueryOptions<BaseResponse<TxSumary>>,
) => {
  return useQuery<BaseResponse<TxSumary>, BaseErrorResponse>({
    queryKey: ['wallet-summary'],
    queryFn: walletService.GetTxSummaryByWallet,
    ...otps,
  });
};

export const useWalletRecentTransactions = (
  params: PaginationParam,
  otps?: QueryOptions<BaseResponseList<RecentTransaction[]>>,
) => {
  return useQuery<BaseResponseList<RecentTransaction[]>, BaseErrorResponse>({
    queryKey: ['wallet-recent-tx', params.limit, params.page],
    queryFn: async () => {
      const res = await walletService.GetWalletRecentTransaction(params);
      return res;
    },
    ...otps,
  });
};
