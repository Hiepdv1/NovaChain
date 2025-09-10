import { GetAddress, GetPublicKeyFromPrivateKey } from '@/lib/db/wallet.store';
import { WalletSignaturePayload } from '../types/wallet';
import { v4 as genuid } from 'uuid';
import { CreateNewTXPayload } from '@/features/tx/types/transaction';

export const buildSignatureWallet = (
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

export const buildSignatureCreateTx = (
  to: string,
  fee: number,
  amount: number,
  timestamp: number,
  pubKey: string,
  message: string,
): CreateNewTXPayload => {
  return {
    data: {
      fee,
      amount,
      to,
      timestamp: timestamp + 3 * 60,
      message,
    },
    pubKey,
    sig: '',
  };
};
