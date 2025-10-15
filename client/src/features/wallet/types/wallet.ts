import { NullTime } from '@/shared/types/api';

export interface WalletConnectData {
  timestamp: number;
  nonce: string;
  publickey: string;
  address: string;
}

export interface WalletConnectPayload {
  data: WalletConnectData;
  sig: string;
}

export interface Wallet {
  ID: string;
  Address: {
    String: string;
    Valid: boolean;
  };
  PublicKey: {
    String: string;
    Valid: boolean;
  };
  PublicKeyHash: string;
  Balance: string;
  CreateAt: NullTime;
  LastLogin: NullTime;
}

export interface WalletSignaturePayload {
  nonce: string;
  publickey: string;
  timestamp: number;
  address: string;
}
