import { StoredWallet } from '@/shared/types/wallet';
import { get, set, del } from 'idb-keyval';
import {
  GetAddress,
  GetPublicKeyFromPrivateKey,
  KeyPair,
} from './wallet.store';
import {
  DecryptedPrivateKeyFromExport,
  EncryptPrivateKeyWithPassword,
} from '../crypto/wallet.crypto';

const WALLET_KEYS_INDEX = 'wallet_pool';

export const AddWalletToStore = async (
  wallet: KeyPair,
  password: string,
): Promise<{
  isFailed: boolean;
  message: string;
}> => {
  try {
    const wallets = await GetWalletPool();
    const address = GetAddress(wallet.publicKey);

    const exists = wallets.some(
      (w) => wallet.publicKey === w.pubkey && w.address === address,
    );

    if (exists) {
      throw new Error('Wallet already exists');
    }

    const encryptedPrivateKey = EncryptPrivateKeyWithPassword(
      wallet.privateKey,
      password,
    );

    const walletData: StoredWallet = {
      address,
      encryptedPrivateKey,
      pubkey: wallet.publicKey,
      createdAt: Date.now(),
    };

    wallets.push(walletData);

    await set(WALLET_KEYS_INDEX, wallets);

    return {
      isFailed: false,
      message: '',
    };
  } catch (err) {
    return {
      isFailed: true,
      message: (err as Error).message,
    };
  }
};

export const GetWalletPool = async (): Promise<StoredWallet[]> => {
  const wallets = await get(WALLET_KEYS_INDEX);
  return wallets || [];
};

export const DelWalletPool = async (): Promise<void> => {
  return await del(WALLET_KEYS_INDEX);
};

export const GetWalletByWalletKey = async (
  encryptedPrivateKey: string,
): Promise<StoredWallet | undefined> => {
  const wallets = await GetWalletPool();
  const privateKey = DecryptedPrivateKeyFromExport(encryptedPrivateKey);

  if (!privateKey) {
    return undefined;
  }
  const publicKey = GetPublicKeyFromPrivateKey(privateKey);
  const address = GetAddress(publicKey);

  return wallets.find((w) => w.pubkey === publicKey && w.address === address);
};

export const DelWalletByWalletKey = async (
  encryptedPrivateKey: string,
): Promise<void> => {
  return await del(encryptedPrivateKey);
};
