import {
  ReAuthModal,
  TransactionDetailModalProps,
  TxVerification,
} from '@/shared/types/modal-type';
import { create } from 'zustand';

export type ModalType =
  | { type: 'reauth'; props: ReAuthModal }
  | {
      type: 'verifyTx';
      props: TxVerification;
    }
  | { type: 'previewTx'; props: TransactionDetailModalProps }
  | { type: null };

export interface ModalActions {
  openModal: (type: ModalType) => void;
  closeModal: () => void;
}

interface ModalState {
  modal: ModalType;
  actions: ModalActions;
}

export const useModalStore = create<ModalState>((set) => ({
  modal: { type: null },
  actions: {
    openModal: (modal: ModalType) => set({ modal }),
    closeModal: () => set({ modal: { type: null } }),
  },
}));
