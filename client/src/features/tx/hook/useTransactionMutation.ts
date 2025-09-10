/* eslint-disable @typescript-eslint/no-explicit-any */
import { BaseErrorResponse, BaseResponse } from '@/shared/types/api';
import { useMutation } from '@tanstack/react-query';
import { CreateNewTXPayload, Transaction } from '../types/transaction';
import transactionService from '../services/transactions.service';

export const useTransactionCreate = () => {
  return useMutation<
    BaseResponse<Transaction>,
    BaseErrorResponse,
    CreateNewTXPayload
  >({
    mutationFn: transactionService.CreateNewTransaction,
  });
};
