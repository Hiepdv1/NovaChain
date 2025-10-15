'use client';

import { useModalStore } from '@/stores/modal-store';
import { Fragment, useEffect, useRef } from 'react';
import AuthModal from '../modals/auth-modal';
import VerifyTransaction from '../modals/tx-verification';
import TransactionDetailModal from '../modals/transaction-detail-modal';
import { useQueryClient } from '@tanstack/react-query';

const ModalProvider = () => {
  const { modal, actions } = useModalStore();
  const queryClient = useQueryClient();
  const prevType = useRef<string | null>(null);

  useEffect(() => {
    if (prevType.current && !modal.type) {
      queryClient.invalidateQueries();
    }
    prevType.current = modal.type;
  }, [modal.type, queryClient]);

  return (
    <Fragment>
      {modal.type === 'reauth' && <AuthModal {...modal.props} {...actions} />}
      {modal.type === 'verifyTx' && (
        <VerifyTransaction {...modal.props} {...actions} />
      )}
      {modal.type === 'previewTx' && (
        <TransactionDetailModal {...modal.props} {...actions} />
      )}
    </Fragment>
  );
};

export default ModalProvider;
