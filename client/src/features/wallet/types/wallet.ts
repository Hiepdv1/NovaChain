import { NullTime } from '@/shared/types/api';
import { NullableString } from '@/shared/types/flag';

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

export interface TxSumary {
  PubKeyHash: string;
  TotalTx: number;
  TotalSent: string;
  TotalReceived: string;
}

export interface RecentTransaction {
  BID: string;
  Type: 'sent' | 'received';
  ID: string;
  TxID: string;
  CreateAt: number;
  Amount: NullableString;
  Fee: NullableString;
  Fromhash: NullableString;
  Tohash: NullableString;
}
