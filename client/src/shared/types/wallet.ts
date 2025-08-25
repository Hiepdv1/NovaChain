export interface StoredWallet {
  address: string;
  pubkey: string;
  encryptedPrivateKey: StoredPrivateKey;
  createdAt: number;
}

export interface StoredPrivateKey {
  cipherText: string;
  authTag: string;
  iv: string;
}
