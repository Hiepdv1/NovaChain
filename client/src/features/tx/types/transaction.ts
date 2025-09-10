export interface TransactionPayload {
  fee: number;
  amount: number;
  to: string;
  message: string;
  timestamp: number;
}

export interface CreateNewTXPayload {
  data: TransactionPayload;
  sig: string;
  pubKey: string;
}

export interface TxInput {
  id: string;
  out: number;
  signature: string | null;
  pubKey: string;
}

export interface TxOutput {
  value: number;
  pubKeyHash: string;
}

export interface Transaction {
  id: string;
  inputs: TxInput[];
  outputs: TxOutput[];
}
