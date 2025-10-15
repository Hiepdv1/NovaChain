import { GetAddress, GetPublicKeyFromPrivateKey } from '@/lib/db/wallet.store';
import { WalletSignaturePayload } from '../types/wallet';
import { v4 as genuid } from 'uuid';
import {
  CreateNewTXPayload,
  SendTransactionData,
  SendTransactionPayload,
} from '@/features/tx/types/transaction';

export const buildWalletSignaturePayload = (
  privateKey: string,
): WalletSignaturePayload => {
  const publickey = GetPublicKeyFromPrivateKey(privateKey);
  const address = GetAddress(publickey);

  return {
    nonce: genuid(),
    publickey,
    timestamp: Math.floor(Date.now() / 1000),
    address,
  };
};

export const buildTransactionSignaturePayload = (
  to: string,
  fee: number,
  amount: number,
  timestamp: number,
  message: string,
  priority: number,
): CreateNewTXPayload => {
  return {
    data: {
      fee,
      amount,
      to,
      timestamp: timestamp + 3 * 60,
      message,
      priority,
    },
    sig: '',
  };
};

export const buildSendTransactionSignaturePayload = (
  payload: SendTransactionData,
): SendTransactionPayload => {
  return {
    data: {
      amount: payload.amount,
      fee: payload.fee,
      message: payload.message,
      receiverAddress: payload.receiverAddress,
      priority: payload.priority,
      transaction: payload.transaction,
    },
    sig: '',
  };
};
