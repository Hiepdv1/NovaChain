import { NullableString } from '@/shared/types/flag';

export interface NetworkOverview {
  Chain: {
    BestHeight: number;
    PerHours: number;
  };
  Hashrate: {
    Value: string;
    Per24H: string;
  };
  Transaction: {
    Total: number;
    AddedToday: number;
  };
  PendingTx: {
    Count: number;
    AddedToday: number;
  };
  ActiveMiners: {
    Count: number;
    Worker: number;
  };
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
}

export interface TransactionItem {
  ID: string;
  TxID: string;
  BID: string;
  CreateAt: number;
  Fromhash: NullableString;
  Tohash: NullableString;
  Fee: NullableString;
  Amount: NullableString;
}

export interface RecentActivityResponse {
  Blocks: BlockItem[];
  Txs: TransactionItem[];
}
