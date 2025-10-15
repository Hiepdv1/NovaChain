/* eslint-disable @typescript-eslint/no-explicit-any */
import { BaseErrorResponse } from '@/shared/types/api';
import { useQuery } from '@tanstack/react-query';
import { Wallet } from '../types/wallet';
import walletService from '../services/wallet.service';
import { QueryOptions } from '@/shared/types/query';

export const useWalletQuery = (opts?: QueryOptions<Wallet>) => {
  return useQuery<Wallet, BaseErrorResponse>({
    queryKey: ['wallet'],
    queryFn: async () => {
      const res = await walletService.GetWallet();
      return res.data;
    },
    ...opts,
  });
};
