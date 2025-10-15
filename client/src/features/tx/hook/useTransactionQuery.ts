import { useQuery } from '@tanstack/react-query';
import transactionService from '../services/transactions.service';
import { BaseErrorResponse, BaseResponseList } from '@/shared/types/api';
import { TransactionPending } from '../types/transaction';
import { QueryOptions } from '@/shared/types/query';

export const useTransactionPending = (
  params: { page?: number; limit?: number } = {},
  opts?: QueryOptions<BaseResponseList<TransactionPending[] | null>>,
) => {
  return useQuery<
    BaseResponseList<TransactionPending[] | null>,
    BaseErrorResponse
  >({
    queryKey: ['txPending'],
    queryFn: async () => {
      const res = await transactionService.GetPendingTxByUser(params);
      return res;
    },
    ...opts,
  });
};
