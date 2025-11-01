import {
  BlockDetail,
  BlockItem,
  GetBlockDetailQuery,
  GetTransactionByBlockSearchQuery,
} from '@/features/block/types/block';
import { handleApiError } from '@/lib/axios/handleErrorApi';
import { http } from '@/lib/axios/http';
import { BaseResponse, BaseResponseList } from '@/shared/types/api';
import { NetworkInfo } from '../types/block';
import { PaginationParam } from '@/shared/types/query';
import { TransactionFull } from '@/features/tx/types/transaction';

class ChainService {
  public async GetNetworkInfo() {
    try {
      const res = await http.get<BaseResponse<NetworkInfo>>('/chain/network');
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }

  public async GetListBlocks(params: PaginationParam) {
    try {
      const res = await http.get<BaseResponseList<BlockItem[]>>(
        `/chain/blocks?page=${params.page}&limit=${params.limit}`,
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }

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

const chainService = new ChainService();

export default chainService;
