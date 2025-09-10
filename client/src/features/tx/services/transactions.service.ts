import { handleApiError } from '@/lib/axios/handleErrorApi';
import { CreateNewTXPayload, Transaction } from '../types/transaction';
import { http } from '@/lib/axios/http';
import { BaseResponse } from '@/shared/types/api';

class TransactionService {
  public async CreateNewTransaction(data: CreateNewTXPayload) {
    try {
      const res = await http.post<BaseResponse<Transaction>>(
        '/txs/__pri/new',
        data,
      );
      return res.data;
    } catch (err) {
      throw handleApiError(err);
    }
  }
}

const transactionService = new TransactionService();

export default transactionService;
