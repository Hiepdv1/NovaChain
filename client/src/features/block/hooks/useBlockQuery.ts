import { useQuery } from '@tanstack/react-query';
import {
  BlockDetail,
  BlockItem,
  GetBlockDetailQuery,
  GetTransactionByBlockSearchQuery,
  NetworkInfo,
} from '../types/block';
import {
  BaseErrorResponse,
  BaseResponse,
  BaseResponseList,
} from '@/shared/types/api';
import { PaginationParam, QueryOptions } from '@/shared/types/query';
import { TransactionFull } from '@/features/tx/types/transaction';
import chainService from '../services/block.service';

export const useBlockDetailQuery = (
  query: GetBlockDetailQuery,
  opts?: QueryOptions<BaseResponse<BlockDetail>>,
) => {
  return useQuery<BaseResponse<BlockDetail>, BaseErrorResponse>({
    queryKey: [query.b_hash, query.page, query.limit],
    queryFn: async () => {
      const res = await chainService.GetBlockDetail(query);
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
      const res = await chainService.GetTransactionByBlockSearch(queries);
      return res;
    },
    ...opts,
  });
};

export const useNetworkInfo = (
  opts?: QueryOptions<BaseResponse<NetworkInfo>>,
) => {
  return useQuery<BaseResponse<NetworkInfo>, BaseErrorResponse>({
    queryKey: ['NetworkInfo'],
    queryFn: chainService.GetNetworkInfo,
    ...opts,
  });
};

export const useListBlocks = (
  queries: PaginationParam,
  opts?: QueryOptions<BaseResponseList<BlockItem[]>>,
) => {
  return useQuery<BaseResponseList<BlockItem[]>, BaseErrorResponse>({
    queryKey: ['blocks', queries.limit, queries.page],
    queryFn: async () => {
      const res = await chainService.GetListBlocks(queries);
      return res;
    },
    ...opts,
  });
};
