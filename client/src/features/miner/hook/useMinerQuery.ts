import { BaseErrorResponse, BaseResponseList } from '@/shared/types/api';
import { useQuery } from '@tanstack/react-query';
import { Miner } from '../types/miner';
import { PaginationParam, QueryOptions } from '@/shared/types/query';
import minerService from '../services/miner.service';

export const useMiners = (
  params: PaginationParam,
  opts?: QueryOptions<BaseResponseList<Miner[]>>,
) => {
  return useQuery<BaseResponseList<Miner[]>, BaseErrorResponse>({
    queryKey: ['miners', params.limit, params.page],
    queryFn: async () => {
      const res = await minerService.GetMiners(params);
      return res;
    },
    ...opts,
  });
};
