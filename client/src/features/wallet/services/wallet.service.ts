import { BaseResponse, BaseResponseList } from '@/shared/types/api';
import {
  RecentTransaction,
  TxSumary,
  Wallet,
  WalletConnectPayload,
} from '../types/wallet';
import { http } from '@/lib/axios/http';
import { handleApiError } from '@/lib/axios/handleErrorApi';
import { EncryptData } from '@/lib/crypto/encode';
import { PaginationParam } from '@/shared/types/query';

class WalletService {
  public async CreateWallet(payload: WalletConnectPayload) {
    try {
      const encode = EncryptData(payload);
      const res = await http.post<BaseResponse<null>>(
        '/wallet/__pub/new',
        encode,
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }

  public async ImportWallet(payload: WalletConnectPayload) {
    try {
      const encode = EncryptData(payload);

      const res = await http.post<BaseResponse<null>>(
        '/wallet/__pub/import',
        encode,
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }

  public async GetWallet() {
    try {
      const res = await http.get<BaseResponse<Wallet>>('/wallet/__pri/me');
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }

  public async Disconnect() {
    try {
      const res = await http.post<BaseResponse<null>>(
        '/wallet/__pri/disconnect',
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }

  public async GetTxSummaryByWallet() {
    try {
      const res = await http.get<BaseResponse<TxSumary>>('/txs/__pri/summary');
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }

  public async GetWalletRecentTransaction(params: PaginationParam) {
    try {
      const res = await http.get<BaseResponseList<RecentTransaction[]>>(
        '/txs/__pri/recent',
        {
          params,
        },
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }
}

const walletService = new WalletService();

export default walletService;
