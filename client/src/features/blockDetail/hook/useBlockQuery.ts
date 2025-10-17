import { useQuery } from '@tanstack/react-query';
import {
  BlockDetail,
  GetBlockDetailQuery,
  GetTransactionByBlockSearchQuery,
} from '../types/block';
import {
  BaseErrorResponse,
  BaseResponse,
  BaseResponseList,
} from '@/shared/types/api';
import blockService from '../services/blockDetail.service';
import { QueryOptions } from '@/shared/types/query';
import { TransactionFull } from '@/features/tx/types/transaction';

export const useBlockDetailQuery = (
  query: GetBlockDetailQuery,
  opts?: QueryOptions<BaseResponse<BlockDetail>>,
) => {
  return useQuery<BaseResponse<BlockDetail>, BaseErrorResponse>({
    queryKey: [query.b_hash, query.page, query.limit],
    queryFn: async () => {
      const res = await blockService.GetBlockDetail(query);
      return res;
    },
    ...opts,
  });
};

export const useSearchTransactionsByBlockQuery = (
  queries: GetTransactionByBlockSearchQuery,
  opts?: QueryOptions<BaseResponseList<TransactionFull[]>>,
) => {
  return useQuery<BaseResponseList<TransactionFull[]>, BaseErrorResponse>({
    queryKey: [queries.b_hash, queries.q, queries.page, queries.limit],
    queryFn: async () => {
      const res = await blockService.GetTransactionByBlockSearch(queries);
      return res;
    },
    ...opts,
  });
};
