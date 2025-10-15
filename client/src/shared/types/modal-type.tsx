/* eslint-disable @typescript-eslint/no-explicit-any */
import {
  ResCreateNewTransaction,
  TransactionPayload,
} from '@/features/tx/types/transaction';
import { StoredWallet } from './wallet';

export interface ReAuthModal {
  title: string;
  des: string;
  notice?: React.ReactNode;
  wallet: {
    address: string;
  };
}

export interface TxVerification {
  onSubmit: (
    wallet: StoredWallet,
    privKey: string,
    data: ResCreateNewTransaction,
  ) => void;
  data: TransactionPayload;
}

export interface TransactionDetailModalProps {
  balance: number;
  from: string;
  to: string;
  fee: number;
  message: string;
  amount: number;
  transaction: ResCreateNewTransaction;
  privateKey: string;
  priority: number;
}

export interface VeriftTxModalProps {
  onSubmit: (
    wallet: StoredWallet,
    privKeyHex: string,
    data: ResCreateNewTransaction,
    ...args: any
  ) => void;
  data: TransactionPayload;
}
