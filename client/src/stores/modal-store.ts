/* eslint-disable @typescript-eslint/no-explicit-any */
import { ReAuthModal } from '@/shared/types/modal-type';
import { create } from 'zustand';

export type ModalType = { type: 'reauth'; props: ReAuthModal } | { type: null };

interface ModalState {
  modal: ModalType;
  actions: {
    openModal: (type: ModalType) => void;
    closeModal: () => void;
  };
}

export const useModalStore = create<ModalState>((set) => ({
  modal: { type: null },
  actions: {
    openModal: (modal: ModalType) => set({ modal }),
    closeModal: () => set({ modal: { type: null } }),
  },
}));
