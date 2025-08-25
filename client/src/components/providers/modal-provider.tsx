'use client';

import { useModalStore } from '@/stores/modal-store';
import { Fragment } from 'react';
import AuthModal from '../modals/auth-modal';

const ModalProvider = ({ children }: { children: React.ReactNode }) => {
  const { modal, actions } = useModalStore();

  if (!modal.type) return children;

  return (
    <Fragment>
      {modal.type === 'reauth' && <AuthModal {...modal.props} {...actions} />}
      {children}
    </Fragment>
  );
};

export default ModalProvider;
