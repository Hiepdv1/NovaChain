/* eslint-disable @typescript-eslint/no-explicit-any */
import {
  Transaction,
  TransactionPayload,
} from '@/features/tx/types/transaction';
import { StoredWallet } from './wallet';
import { ModalType } from '@/stores/modal-store';

export interface ReAuthModal {
  title: string;
  des: string;
  notice?: React.ReactNode;
  wallet: {
    address: string;
  };
}

export interface TxVerification {
  onSubmit: (wallet: StoredWallet, privKey: string, data: Transaction) => void;
  data: TransactionPayload;
}

export interface TransactionDetailModalProps {
  balance: number;
  from: string;
  to: string;
  fee: number;
  message: string;
  amount: number;
  transactions: Transaction;
}

export interface VeriftTxModalProps {
  onSubmit: (
    wallet: StoredWallet,
    privKeyHex: string,
    data: Transaction,
    ...args: any
  ) => void;
  data: TransactionPayload;
}
