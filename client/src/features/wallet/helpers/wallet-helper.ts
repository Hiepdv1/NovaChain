import { GetAddress, GetPublicKeyFromPrivateKey } from '@/lib/db/wallet.store';
import { WalletSignaturePayload } from '../types/wallet';
import { v4 as genuid } from 'uuid';

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
