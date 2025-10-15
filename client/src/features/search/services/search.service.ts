import { handleApiError } from '@/lib/axios/handleErrorApi';
import { SearchItem, SearchQuery } from '../types/search';
import { http } from '@/lib/axios/http';
import { BaseResponseList } from '@/shared/types/api';

class SearchService {
  public async GetSearchResult(query: SearchQuery) {
    try {
      const res = await http.get<BaseResponseList<SearchItem[] | null>>(
        `/chain/search?search_query=${query.search_query}&page=${query.page}&${query.limit}`,
      );

      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }
}

const searchService = new SearchService();

export default searchService;
