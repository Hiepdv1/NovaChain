import { handleApiError } from '@/lib/axios/handleErrorApi';
import { http } from '@/lib/axios/http';
import {
  BlockDetail,
  GetBlockDetailQuery,
  GetTransactionByBlockSearchQuery,
} from '../types/block';
import { BaseResponse, BaseResponseList } from '@/shared/types/api';
import { TransactionFull } from '@/features/tx/types/transaction';

class BlockService {
  public async GetBlockDetail(query: GetBlockDetailQuery) {
    try {
      const res = await http.get<BaseResponse<BlockDetail>>(
        `/chain/blocks/${query.b_hash}?page=${query.page}&limit=${query.limit}`,
      );

      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }

  public async GetTransactionByBlockSearch(
    query: GetTransactionByBlockSearchQuery,
  ) {
    try {
      const res = await http.get<BaseResponseList<TransactionFull[]>>(
        `/txs/__pub/search?b_hash=${query.b_hash}&q=${query.q}&page=${query.page}&limit=${query.limit}`,
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }
}

const blockService = new BlockService();

export default blockService;
