import { useQuery } from '@tanstack/react-query';
import transactionService from '../services/transactions.service';
import { BaseErrorResponse, BaseResponseList } from '@/shared/types/api';
import { TransactionItem, TransactionPending } from '../types/transaction';
import { PaginationParam, QueryOptions } from '@/shared/types/query';

export const useTransactionPending = (
  params: PaginationParam,
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

export const useListTransactions = (
  params: PaginationParam,
  opts?: QueryOptions<BaseResponseList<TransactionItem[]>>,
) => {
  return useQuery<BaseResponseList<TransactionItem[]>, BaseErrorResponse>({
    queryKey: ['txs', params.limit, params.page],
    queryFn: async () => {
      const res = await transactionService.GetListTransactions(params);
      return res;
    },
    ...opts,
  });
};

export const usePendingTransactions = (
  params: PaginationParam,
  opts?: QueryOptions<BaseResponseList<TransactionPending[]>>,
) => {
  return useQuery<BaseResponseList<TransactionPending[]>, BaseErrorResponse>({
    queryKey: ['txPendings', params.limit, params.page],
    queryFn: async () => {
      const res = await transactionService.GetPendingTxs(params);
      return res;
    },
    ...opts,
  });
};
