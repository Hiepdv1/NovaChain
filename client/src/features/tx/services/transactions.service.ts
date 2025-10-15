import { handleApiError } from '@/lib/axios/handleErrorApi';
import {
  CreateNewTXPayload,
  ResCreateNewTransaction,
  SendTransactionPayload,
  TransactionPending,
} from '../types/transaction';
import { http } from '@/lib/axios/http';
import { BaseResponse, BaseResponseList } from '@/shared/types/api';
import { DecryptData, EncryptData } from '@/lib/crypto/encode';

class TransactionService {
  public async CreateNewTransaction(data: CreateNewTXPayload) {
    try {
      const encode = EncryptData(data);

      const res = await http.post<BaseResponse<string>>(
        '/txs/__pri/new',
        encode,
      );
      const txEncrypted = res.data;
      const decrypted = DecryptData<ResCreateNewTransaction>(txEncrypted.data);

      const baseRes: BaseResponse<ResCreateNewTransaction> = {
        ...txEncrypted,
        data: decrypted,
      };

      return baseRes;
    } catch (err) {
      throw handleApiError(err);
    }
  }

  public async SendTransaction(data: SendTransactionPayload) {
    try {
      const encode = EncryptData(data);
      const res = await http.post<BaseResponse<null>>(
        '/txs/__pri/send',
        encode,
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }

  public async GetPendingTxByUser(params: { page?: number; limit?: number }) {
    try {
      const res = await http.get<BaseResponseList<TransactionPending[] | null>>(
        '/txs/__pri/pending',
        {
          params: {
            limit: params.limit || 1,
            page: params.page || 1,
          },
        },
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }
}

const transactionService = new TransactionService();

export default transactionService;
