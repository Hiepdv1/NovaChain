import { NullTime } from '@/shared/types/api';
import { QueryKey, UseQueryOptions } from '@tanstack/react-query';

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
  Address: string;
  PublicKey: string;
  PublicKeyHash: string;
  Balance: string;
  CreateAt: NullTime;
  LastLogin: NullTime;
}

export type WalletQueryOptions = Omit<
  UseQueryOptions<Wallet, Error, Wallet, QueryKey>,
  'queryKey' | 'queryFn'
>;

export interface WalletSignaturePayload {
  nonce: string;
  publickey: string;
  timestamp: number;
  address: string;
}
