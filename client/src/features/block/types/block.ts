import { TransactionFull } from '@/features/tx/types/transaction';
import { PaginationMeta } from '@/shared/types/api';
import { NullableString } from '@/shared/types/flag';
import { PaginationParam } from '@/shared/types/query';

export interface GetBlockDetailQuery extends PaginationParam {
  b_hash: string;
}
export interface BlockDetail {
  ID: string;
  BID: string;
  PrevHash: NullableString;
  Nonce: number;
  Height: number;
  MerkeleRoot: string;
  Nbits: number;
  TxCount: number;
  NchainWork: string;
  Size: number;
  Timestamp: number;
  TotalFee: string;
  Difficulty: number;
  Miner: string;
  Transactions: {
    Data: TransactionFull[];
    Meta: PaginationMeta;
  };
}

export interface GetTransactionByBlockSearchQuery extends PaginationParam {
  b_hash: string;
  q: string;
}

export interface NetworkInfo {
  LastBlock: number;
  Hashrate: string;
  AvgBlockTime: number;
  AvgDifficulty: number;
  SyncStatus: string;
  NetworkHealth: string;
  TxPending: number;
}

export interface BlockItem {
  ID: string;
  BID: string;
  PrevHash: NullableString;
  Nonce: number;
  Height: number;
  MerkleRoot: string;
  Nbits: number;
  TxCount: number;
  NchainWork: string;
  Size: number;
  Timestamp: number;
  PubKeyHash: string;
  Value: string;
}
