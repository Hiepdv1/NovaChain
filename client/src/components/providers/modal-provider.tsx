'use client';

import { useModalStore } from '@/stores/modal-store';
import { Fragment } from 'react';
import AuthModal from '../modals/auth-modal';
import VerifyTransaction from '../modals/tx-verification';
import TransactionDetailModal from '../modals/transaction-detail-modal';

const ModalProvider = () => {
  const { modal, actions } = useModalStore();

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
