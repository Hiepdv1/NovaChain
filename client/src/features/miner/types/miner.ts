export interface Miner {
  MinerPubkey: string;
  MinedBlocks: number;
  FirstMinedAt: number;
  LastMinedAt: number;
  TotalBlocksNetwork: number;
  NetworkSharePercent: string;
}
