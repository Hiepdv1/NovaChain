/* eslint-disable @typescript-eslint/no-explicit-any */
import { BaseErrorResponse, BaseResponse } from '@/shared/types/api';
import { useMutation } from '@tanstack/react-query';
import {
  CreateNewTXPayload,
  ResCreateNewTransaction,
  SendTransactionPayload,
} from '../types/transaction';
import transactionService from '../services/transactions.service';

export const useTransactionCreate = () => {
  return useMutation<
    BaseResponse<ResCreateNewTransaction>,
    BaseErrorResponse,
    CreateNewTXPayload
  >({
    mutationFn: transactionService.CreateNewTransaction,
  });
};

export const useSendTransactionMutation = () => {
  return useMutation<
    BaseResponse<null>,
    BaseErrorResponse,
    SendTransactionPayload
  >({
    mutationFn: transactionService.SendTransaction,
  });
};
