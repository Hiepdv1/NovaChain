'use client';

import { useCallback, useState } from 'react';
import Welcome from './components/welcome';
import CreateWallet from './components/create-wallet';
import ImportWallet from './components/import-wallet';

export type ModalName = 'Welcome' | 'CreateWallet' | 'ImportWallet';

const WalletPage = () => {
  const [showModal, setShowModal] = useState({
    isShowWelcome: true,
    isShowCreateWallet: false,
    isShowImportWallet: false,
  });

  const onSwitchModal = useCallback((modelName: ModalName) => {
    switch (modelName) {
      case 'Welcome':
        return setShowModal(() => {
          return {
            isShowWelcome: true,
            isShowCreateWallet: false,
            isShowImportWallet: false,
          };
        });
      case 'CreateWallet':
        return setShowModal(() => {
          return {
            isShowWelcome: false,
            isShowCreateWallet: true,
            isShowImportWallet: false,
          };
        });
      case 'ImportWallet':
        return setShowModal(() => {
          return {
            isShowWelcome: false,
            isShowCreateWallet: false,
            isShowImportWallet: true,
          };
        });
    }
  }, []);

  return (
    <div className="w-full max-w-md relative z-10">
      {showModal.isShowWelcome && <Welcome onSwitchModal={onSwitchModal} />}
      {showModal.isShowCreateWallet && (
        <CreateWallet onSwitchModal={onSwitchModal} />
      )}
      {showModal.isShowImportWallet && (
        <ImportWallet onSwitchModal={onSwitchModal} />
      )}
    </div>
  );
};

export default WalletPage;
