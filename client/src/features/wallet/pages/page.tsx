'use client';

import { Fragment, useCallback, useState } from 'react';
import Welcome from '../components/welcome';
import CreateWallet from '../components/create-wallet';
import ImportWallet from '../components/import-wallet';
import LoadingOverlay, {
  LoadingOverlayProps,
  LoadingStep,
} from '../components/loading-overlay';

export type ModalName = 'Welcome' | 'CreateWallet' | 'ImportWallet';

const createWalletInitSteps: LoadingStep[] = [
  {
    active: false,
    completed: false,
    title: 'Initializing secure connection...',
  },
  {
    active: false,
    completed: false,
    title: 'Generating cryptographic keys...',
  },
  { active: false, completed: false, title: 'Encrypting wallet data...' },
  { active: false, completed: false, title: 'Finalizing wallet creation...' },
];

const loadingCreateWallet = {
  isLoading: false,
  modalData: {
    title: 'Create Your Wallet',
    des: 'Initializing secure wallet generation...',
    loadingSteps: createWalletInitSteps,
  },
};

const importWalletInitSteps: LoadingStep[] = [
  {
    active: false,
    completed: false,
    title: 'Initializing secure connection...',
  },
  {
    active: false,
    completed: false,
    title: 'Validating wallet credentials...',
  },
  { active: false, completed: false, title: 'Decrypting wallet data...' },
  {
    active: false,
    completed: false,
    title: 'Importing wallet and synchronizing...',
  },
];

const loadingImportWallet = {
  isLoading: false,
  modalData: {
    title: 'Import Existing Wallet',
    des: 'Preparing to import and synchronize wallet...',
    loadingSteps: importWalletInitSteps,
  },
};

const WalletPage = () => {
  const [showModal, setShowModal] = useState({
    isShowWelcome: true,
    isShowCreateWallet: false,
    isShowImportWallet: false,
  });

  const [createWalletLoading, setCreateWalletLoading] = useState<{
    isLoading: boolean;
    modalData: LoadingOverlayProps;
  }>(loadingCreateWallet);

  const [importWalletLoading, setImportWalletLoading] = useState<{
    isLoading: boolean;
    modalData: LoadingOverlayProps;
  }>(loadingImportWallet);

  const updateCreateWalletStep = useCallback(
    (stepNumber = 1, isLoading: boolean) => {
      setCreateWalletLoading((prev) => {
        return {
          ...prev,
          isLoading: isLoading,
          modalData: {
            ...prev.modalData,
            loadingSteps: prev.modalData.loadingSteps.map((step, ix) => {
              if (ix === stepNumber - 1) {
                return { ...step, active: true, completed: false };
              } else if (ix < stepNumber - 1) {
                return { ...step, active: false, completed: true };
              }
              return { ...step, active: false, completed: false };
            }),
          },
        };
      });
    },
    [],
  );

  const updateImportWalletStep = useCallback(
    (isLoading = false, stepNumber = 1) => {
      setImportWalletLoading((prev) => {
        return {
          ...prev,
          isLoading,
          modalData: {
            ...prev.modalData,
            loadingSteps: prev.modalData.loadingSteps.map((step, ix) => {
              if (ix === stepNumber - 1) {
                return { ...step, active: true, completed: false };
              } else if (ix < stepNumber - 1) {
                return { ...step, active: false, completed: true };
              }
              return { ...step, active: false, completed: false };
            }),
          },
        };
      });
    },
    [],
  );

  const showModalByName = useCallback((modelName: ModalName) => {
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
    <Fragment>
      <div className="w-full max-w-md relative z-10">
        {showModal.isShowWelcome && <Welcome onSwitchModal={showModalByName} />}
        {showModal.isShowCreateWallet && (
          <CreateWallet
            showModalByName={showModalByName}
            onStepUpdate={updateCreateWalletStep}
          />
        )}
        {showModal.isShowImportWallet && (
          <ImportWallet
            onStepUpdate={updateImportWalletStep}
            onSwitchModal={showModalByName}
          />
        )}
      </div>
      {createWalletLoading.isLoading && (
        <LoadingOverlay
          loadingSteps={createWalletLoading.modalData.loadingSteps}
          title={createWalletLoading.modalData.title}
          des={createWalletLoading.modalData.des}
        />
      )}

      {importWalletLoading.isLoading && (
        <LoadingOverlay
          loadingSteps={importWalletLoading.modalData.loadingSteps}
          title={importWalletLoading.modalData.title}
          des={importWalletLoading.modalData.des}
        />
      )}
    </Fragment>
  );
};

export default WalletPage;
