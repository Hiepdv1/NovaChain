import { PaginationParam } from '@/shared/types/query';

export interface SearchBlockItem {
  Type: 'block';
  Keyword: string;
  Data: {
    size: number;
    miner: string;
    height: number;
    tx_count: number;
    timestamp: number;
  };
}

export interface SearchTxItem {
  Type: 'transaction';
  Keyword: string;
  Data: {
    to: string;
    fee: number;
    from: string;
    amount: number;
    timestamp: number;
  };
}

export type SearchItem = SearchBlockItem | SearchTxItem;

export interface SearchQuery extends PaginationParam {
  search_query: string;
}
