import { BlockItem } from '@/features/block/types/block';
import { TransactionItem } from '@/features/tx/types/transaction';

export interface NetworkOverview {
  Chain: {
    BestHeight: number;
    PerHours: number;
  };
  Hashrate: {
    Value: string;
    ChangeRate: string;
    Trend: 'increase' | 'decrease' | 'stable';
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

export interface RecentActivityResponse {
  Blocks: BlockItem[];
  Txs: TransactionItem[];
}
