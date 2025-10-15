'use client';

import { Fragment, memo, MouseEvent, useCallback, useState } from 'react';
import { ModalName } from '../pages/page';
import Button from '@/components/button';
import CreatePassword from './create-password';
import { toast } from '@/components/globalToaster';
import { useRouter } from 'next/navigation';

import { useWalletConnect } from '../hook/useWalletMutation';
import { v4 as uuid } from 'uuid';
import { CreateWallet, GetAddress } from '@/lib/db/wallet.store';
import { AddWalletToStore, DelWalletByWalletKey } from '@/lib/db/wallet.index';
import {
  EncryptPrivateKeyForExport,
  SignPayload,
} from '@/lib/crypto/wallet.crypto';
import NoticeBox from './notice-box';
import { WalletSignaturePayload } from '../types/wallet';
import useWalletContext from '@/components/providers/wallet-provider';

interface CreateWallet {
  showModalByName(modelName: ModalName): void;
  onStepUpdate: (stepNumber: number, isLoading: boolean) => void;
}

type Wallet = {
  address: string;
  privKey: string;
  pubKey: string;
};

const CreateWalletForm = ({ showModalByName, onStepUpdate }: CreateWallet) => {
  const [status, setStatus] = useState<{
    isShowCreatePassWord: boolean;
    isShowGenerateWallet: boolean;
  }>({
    isShowCreatePassWord: true,
    isShowGenerateWallet: false,
  });

  const [keyPair, setKeyPair] = useState<Wallet>({
    address: '',
    privKey: '',
    pubKey: '',
  });

  const router = useRouter();
  const { refetch } = useWalletContext();
  const walletConnect = useWalletConnect();

  const onBackCreatePassword = useCallback(() => {
    showModalByName('Welcome');
  }, [showModalByName]);

  const onContinue = useCallback(
    async (password: string) => {
      onStepUpdate(1, true);
      let wallet = CreateWallet();
      while (
        wallet.publicKey.length !== 128 ||
        wallet.privateKey.length !== 64
      ) {
        wallet = CreateWallet();
      }

      const privKeyEncrypted = EncryptPrivateKeyForExport(wallet.privateKey);

      onStepUpdate(2, true);

      const { isFailed, message } = await AddWalletToStore(wallet, password);

      onStepUpdate(3, true);

      if (isFailed) {
        toast.error(message);
        return;
      }

      const address = GetAddress(wallet.publicKey);

      const data: WalletSignaturePayload = {
        nonce: uuid(),
        publickey: wallet.publicKey,
        timestamp: Math.floor(Date.now() / 1000),
        address,
      };

      const sig = SignPayload(wallet.privateKey, data);

      onStepUpdate(4, true);

      walletConnect.mutate(
        {
          data,
          sig,
        },
        {
          onError: async () => {
            await DelWalletByWalletKey(privKeyEncrypted);
            toast.error('Create wallet failed. Please try again.');
            onStepUpdate(1, false);
          },
          onSuccess: () => {
            setKeyPair({
              privKey: privKeyEncrypted,
              pubKey: wallet.publicKey,
              address,
            });

            setStatus(() => {
              return {
                isShowCreatePassWord: false,
                isShowGenerateWallet: true,
              };
            });

            toast.success('Wallet created successfully 🚀');

            onStepUpdate(1, false);
          },
        },
      );
    },

    [walletConnect, onStepUpdate],
  );

  const onCopyToClipboard = useCallback((e: MouseEvent<HTMLButtonElement>) => {
    const button = e.currentTarget;
    const parent = button.parentElement;
    const code = parent?.childNodes[0].textContent || '';
    const originalText = button.textContent;

    navigator.clipboard
      .writeText(code)
      .then(() => {
        button.textContent = 'Copied!';
        button.style.background = 'rgba(34, 197, 94, 0.3)';

        setTimeout(() => {
          button.textContent = originalText;
          button.style.background = '';
        }, 2000);

        toast.success('Copied to clipboard!');
      })
      .catch(() => {
        toast.error('Failed to copy');
      });
  }, []);

  const onBackToHome = useCallback(async () => {
    if (!refetch) return;
    await refetch();
    onStepUpdate(1, false);
    router.push('/');
    router.refresh();
  }, [router, refetch, onStepUpdate]);

  const onAccessWallet = useCallback(async () => {
    if (!refetch) return;
    await refetch();
    onStepUpdate(1, false);
    router.push('/wallet/me');
    router.refresh();
  }, [router, refetch, onStepUpdate]);

  return (
    <Fragment>
      {/* Step 1: Create Password */}
      {status.isShowCreatePassWord && (
        <div className="glass-card rounded-3xl p-8">
          <div
            style={{
              animationDelay: '100ms',
            }}
            className="quantum-steps animate-cascase-fade mb-12"
          >
            <div className="quantum-step before:!left-[92%] before:w-16 before:h-0.5 active">
              <div className="quantum-circle w-12 h-12 active animate-quantum-active-grow">
                <span className="text-white font-bold text-sm">1</span>
              </div>
              <p className="text-white text-xs font-semibold mt-4 text-center">
                Create Password
              </p>
            </div>

            <div className="quantum-step before:!bg-transparent">
              <div className="quantum-circle w-12 h-12">
                <span className="text-white font-bold text-sm">2</span>
              </div>
              <p className="text-white text-xs font-semibold mt-4 text-center">
                Generate Wallet
              </p>
            </div>
          </div>

          <CreatePassword
            onContinue={onContinue}
            onBack={onBackCreatePassword}
          />
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
                    {keyPair.address}
                  </code>
                  <Button
                    variant="glass"
                    size="sm"
                    onClick={onCopyToClipboard}
                    className="shrink-0 w-fit text-gray-700 text-[10px] rounded-lg hover:bg-white"
                  >
                    Copy
                  </Button>
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
                    {keyPair.privKey}
                  </code>
                  <Button
                    variant="glass"
                    size="sm"
                    onClick={onCopyToClipboard}
                    className="shrink-0 w-fit text-gray-700 text-[10px] rounded-lg hover:bg-white"
                  >
                    Copy
                  </Button>
                </div>
              </div>
            </div>

            <NoticeBox
              description="Your private key is your digital signature. Store it safely
                    offline and never share it with anyone."
              icon={
                <svg
                  className="w-8 h-8 text-orange-500 mt-0.5"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path d="M12,2L1,21H23M12,6L19.53,19H4.47M11,10V14H13V10M11,16V18H13V16"></path>
                </svg>
              }
              title="Security Reminder"
              variant="error"
            />

            <div className="mt-8 space-y-3">
              <Button onClick={onAccessWallet} variant="glass" size="md">
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

              <Button onClick={onBackToHome} variant="glass" size="md">
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

export default memo(CreateWalletForm);
