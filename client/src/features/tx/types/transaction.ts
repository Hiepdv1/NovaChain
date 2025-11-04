import { NullableString } from '@/shared/types/flag';

export interface TransactionPayload {
  fee: number;
  amount: number;
  to: string;
  timestamp: number;
  priority: number;
}

export interface CreateNewTXPayload {
  data: TransactionPayload;
  sig: string;
}

export interface TxInput {
  ID: string;
  Out: number;
  Signature: string | null;
  PubKey: string;
}

export interface TxOutput {
  Value: number;
  PubKeyHash: string;
}

export interface Transaction {
  ID: string;
  Inputs: TxInput[];
  Outputs: TxOutput[];
}

export interface SendTransactionData {
  amount: number;
  fee: number;
  transaction: Transaction;
  message: string;
  receiverAddress: string;
  priority: number;
}

export interface SendTransactionPayload {
  data: SendTransactionData;
  sig: string;
}

export interface TxInputWithDataToSign {
  id: string;
  out: number;
  signature: string | null;
  pubKey: string;
  dataToSign: string;
}

export interface UTXO {
  ID: string;
  TxID: {
    String: string;
    Valid: boolean;
  };
  OutputIndex: number;
  Value: string;
  PubKeyHash: string;
  BlockID: string;
}

export interface ResCreateNewTransaction {
  id: string;
  inputs: TxInputWithDataToSign[];
  outputs: TxOutput[];
}

export interface TransactionPending {
  ID: string;
  TxID: string;
  Address: string;
  ReceiverAddress: string;
  Amount: string;
  Fee: string;
  Status: 'pending' | 'failed';
  Priority: {
    Int32: number;
    Valid: boolean;
  };
  Message: {
    String: string;
    Valid: boolean;
  };
  CreatedAt: {
    Time: string;
    Valid: boolean;
  };
  UpdatedAt: {
    Time: string;
    Valid: boolean;
  };
}

export interface TransactionFull {
  ID: string;
  TxID: string;
  BID: string;
  CreateAt: number;
  Amount: NullableString;
  Fee: NullableString;
  Fromhash: NullableString;
  Tohash: NullableString;
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
