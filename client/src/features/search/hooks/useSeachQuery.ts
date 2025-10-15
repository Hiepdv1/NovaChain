import { BaseErrorResponse, BaseResponseList } from '@/shared/types/api';
import { QueryOptions } from '@/shared/types/query';
import { SearchItem, SearchQuery } from '../types/search';
import { useQuery } from '@tanstack/react-query';
import searchService from '../services/search.service';

export const useSearchQuery = (
  query: SearchQuery,
  opts?: QueryOptions<BaseResponseList<SearchItem[] | null>>,
) => {
  return useQuery<BaseResponseList<SearchItem[] | null>, BaseErrorResponse>({
    queryKey: [query.search_query, query.page, query.limit],
    queryFn: async () => {
      const res = await searchService.GetSearchResult(query);
      return res;
    },
    ...opts,
  });
};
