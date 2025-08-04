import { Fragment, memo, useCallback, useState } from 'react';
import { ModalName } from '../page';
import Button from '@/components/button';
import CreatePassword from './create-password';

interface CreateWallet {
  onSwitchModal(modelName: ModalName): void;
}

const CreateWallet = ({ onSwitchModal }: CreateWallet) => {
  const [status, setStatus] = useState<{
    isShowCreatePassWord: boolean;
    isShowGenerateWallet: boolean;
  }>({
    isShowCreatePassWord: true,
    isShowGenerateWallet: false,
  });

  const onShowWelcomeModal = () => {
    onSwitchModal('Welcome');
  };

  const onBackCreatePassword = useCallback(() => {
    onSwitchModal('Welcome');
  }, [onSwitchModal]);

  return (
    <Fragment>
      {/* Step 1: Create Password */}
      {status.isShowCreatePassWord && (
        <div className="glass-card rounded-3xl  p-8">
          <div
            style={{
              animationDelay: '100ms',
            }}
            className="quantum-steps animate-cascase-fade mb-12"
          >
            <div className="quantum-step before:h-0.5 active">
              <div className="quantum-circle w-12 h-12 active animate-quantum-active-grow">
                <span className="text-white font-bold text-sm">1</span>
              </div>
              <p className="text-white text-xs font-semibold mt-4 text-center">
                Create Password
              </p>
            </div>

            <div className="quantum-step">
              <div className="quantum-circle before:content-none w-12 h-12">
                <span className="text-white font-bold text-sm">2</span>
              </div>
              <p className="text-white text-xs font-semibold mt-4 text-center">
                Generate Wallet
              </p>
            </div>
          </div>

          <CreatePassword onBack={onBackCreatePassword} />
        </div>
      )}

      {/* Step 2: Generate Wallet */}
      {status.isShowGenerateWallet && (
        <div className="glass-card rounded-3xl p-8 animate-slide-up">
          <div className="text-center mb-8">
            <div className="w-20 h-20 mb-4 mx-auto rounded-full success-button flex items-center justify-center animate-bounce-in">
              <svg
                className="w-10 h-10 text-white"
                fill="currentColor"
                viewBox="0 0 24 24"
              >
                <path d="M9,20.42L2.79,14.21L5.62,11.38L9,14.77L18.88,4.88L21.71,7.71L9,20.42Z"></path>
              </svg>
            </div>
            <h2 className="text-2xl font-bold text-white mb-2 text-shadow-2xs animate-floating-text">
              Wallet Created!
            </h2>
            <p
              style={{
                animationDelay: '0.5s',
              }}
              className="text-white text-[16px] opacity-0 whitespace-nowrap overflow-hidden border-r-2 border-[rgba(255,255,255,0.7)] animate-writer"
            >
              Your secure digital wallet is ready to use
            </p>
          </div>

          <div className="space-y-6">
            <div>
              <div className="flex items-center justify-between mb-3">
                <label className="text-white text-sm font-medium">
                  Public Address
                </label>
                <span className="text-[10px] text-green-600 bg-green-100 px-2 py-1 rounded-full font-medium">
                  SHAREABLE
                </span>
              </div>
              <div className="key-display rounded-xl p-4 rainbow-shimmer">
                <div className="flex items-center justify-between relative z-10">
                  <code className="text-xs pr-3 text-gray-700 font-mono break-all">
                    0xc233d8e54ced2a161aee177530c961ecc50963b8
                  </code>
                  <button className="shrink-0 text-[10px] cursor-pointer glass-button text-gray-700 px-4 py-2 rounded-lg font-medium hover:bg-white">
                    Copy
                  </button>
                </div>
              </div>
            </div>

            <div>
              <div className="flex items-center justify-between mb-3">
                <label className="text-white text-sm font-medium">
                  Private Key
                </label>
                <span className="text-[10px] text-red-600 bg-red-200 px-2 py-1 rounded-full font-medium">
                  KEEP SECRET
                </span>
              </div>

              <div className="key-display rounded-xl p-4 shimmer">
                <div className="flex items-center justify-between relative z-10">
                  <code className="text-xs pr-3 text-gray-700 font-mono break-all">
                    0x13c80f1e41c97af74bd3bdd1fd5b5cdd71e5452702af0ed4628633e1d2793578
                  </code>
                  <button className="shrink-0 text-[10px] cursor-pointer glass-button text-gray-700 px-4 py-2 rounded-lg font-medium hover:bg-white">
                    Copy
                  </button>
                </div>
              </div>
            </div>

            <div className="glass-card rounded-xl p-4 border-l-4 border-orange-600">
              <div className="flex items-start space-x-3">
                <svg
                  className="w-6 h-6 text-orange-500 mt-0.5"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path d="M12,2L1,21H23M12,6L19.53,19H4.47M11,10V14H13V10M11,16V18H13V16"></path>
                </svg>
                <div>
                  <p className="text-white font-semibold text-sm">
                    Security Reminder
                  </p>
                  <p className="text-white opacity-80 text-xs mt-1 leading-relaxed">
                    Your private key is your digital signature. Store it safely
                    offline and never share it with anyone.
                  </p>
                </div>
              </div>
            </div>

            <div className="mt-8 space-y-3">
              <Button variant="glass" size="md">
                <div className="flex items-center justify-center text-[14px]">
                  <span>Access wallet</span>
                  <svg
                    className="w-5 h-5 ml-2"
                    fill="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path d="M4,11V13H16L10.5,18.5L11.92,19.92L19.84,12L11.92,4.08L10.5,5.5L16,11H4Z"></path>
                  </svg>
                </div>
              </Button>

              <Button onClick={onShowWelcomeModal} variant="glass" size="md">
                <div className="flex items-center justify-center text-[14px]">
                  <svg
                    className="w-5 h-5 mr-2"
                    fill="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path d="M10,20V14H14V20H19V12H22L12,3L2,12H5V20H10Z"></path>
                  </svg>
                  <span>Back to home</span>
                </div>
              </Button>
            </div>
          </div>
        </div>
      )}
    </Fragment>
  );
};

export default memo(CreateWallet);
