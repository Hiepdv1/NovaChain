'use client';

import {
  ChangeEvent,
  Fragment,
  MouseEvent,
  useCallback,
  useEffect,
  useRef,
  useState,
} from 'react';
import { ModalName } from '../pages/page';
import CreatePassword from './create-password';
import Button from '@/components/button';
import { toast } from '@/components/globalToaster';
import {
  AddWalletToStore,
  DelWalletByWalletKey,
  GetWalletByWalletKey,
} from '@/lib/db/wallet.index';
import {
  DecryptedPrivateKeyFromExport,
  SignPayload,
} from '@/lib/crypto/wallet.crypto';
import { useWalletImport } from '../hook/useWalletQuery';
import { WalletConnectData } from '../types/wallet';
import { GetAddress, GetPublicKeyFromPrivateKey } from '@/lib/db/wallet.store';
import { v4 as uuid } from 'uuid';
import { useRouter } from 'next/navigation';
import useWalletContext from '@/components/providers/wallet-provider';

interface ImportWallet {
  onSwitchModal(modalName: ModalName): void;
  onStepUpdate(isLoading?: boolean, stepNumber?: number): void;
}

type MapRef = {
  continueBtn: HTMLButtonElement | null;
  validationDiv: HTMLDivElement | null;
  textArea: HTMLTextAreaElement | null;
  privateKey: {
    value: string;
    encode: string;
  };
};

const ImportWallet = ({ onSwitchModal, onStepUpdate }: ImportWallet) => {
  const [status, setStatus] = useState<{
    isShowImportWallet: boolean;
    isShowCreatePassWord: boolean;
    isCompleted: boolean;
    currentStep: number;
  }>({
    isShowImportWallet: true,
    isShowCreatePassWord: false,
    isCompleted: false,
    currentStep: 1,
  });

  const { refetch } = useWalletContext();
  const router = useRouter();

  const refs = useRef<MapRef>({
    continueBtn: null,
    validationDiv: null,
    textArea: null,
    privateKey: {
      value: '',
      encode: '',
    },
  });

  const [importWallet, setImportWallet] = useState<{
    address: string;
    pubkey: string;
  }>({
    address: '',
    pubkey: '',
  });

  const walletImport = useWalletImport();

  const onContinue = useCallback(async () => {
    const { textArea } = refs.current;
    if (!textArea) {
      toast.error('Not available yet. Please try again later.');
      return;
    }

    const privKeyEncrypted = textArea.value;

    const exists = await GetWalletByWalletKey(privKeyEncrypted);

    if (exists) {
      toast.error(
        'This wallet is imported already. Please choose another one.',
      );
      return;
    }

    const privKey = DecryptedPrivateKeyFromExport(privKeyEncrypted);
    if (!privKey) {
      toast.error('Invalid private key format. Please try again.');
      return;
    }

    refs.current.privateKey.value = privKey;
    refs.current.privateKey.encode = privKeyEncrypted;

    const pubkey = GetPublicKeyFromPrivateKey(privKey);
    const address = GetAddress(pubkey);

    setImportWallet({
      pubkey,
      address,
    });

    setStatus(() => ({
      isCompleted: false,
      isShowImportWallet: false,
      isShowCreatePassWord: true,
      currentStep: 2,
    }));

    toast.success('🎉 Wallet imported successfully');
  }, []);

  const onBackCreatePassword = useCallback(() => {
    setStatus(() => ({
      isCompleted: false,
      isShowImportWallet: true,
      isShowCreatePassWord: false,
      currentStep: 1,
    }));
  }, []);

  const onContinueCreatePassword = useCallback(
    async (password: string) => {
      onStepUpdate(true, 1);
      const { privateKey } = refs.current;

      onStepUpdate(true, 2);
      if (privateKey.value === '' || privateKey.encode === '') {
        toast.error('Something wrong, please try again later.');
        return;
      }

      const publicKey = GetPublicKeyFromPrivateKey(privateKey.value);
      const address = GetAddress(publicKey);

      onStepUpdate(true, 3);

      const data: WalletConnectData = {
        nonce: uuid(),
        publickey: publicKey,
        timestamp: Math.floor(Date.now() / 1000),
        address,
      };

      const { isFailed, message } = await AddWalletToStore(
        {
          privateKey: privateKey.value,
          publicKey,
        },
        password,
      );

      if (isFailed) {
        toast.error(message);
        return;
      }

      const sig = SignPayload(privateKey.value, data);
      onStepUpdate(true, 4);
      walletImport.mutate(
        {
          data,
          sig,
        },
        {
          onError: async () => {
            await DelWalletByWalletKey(privateKey.encode);
            onStepUpdate(false);
          },
          onSuccess: () => {
            onStepUpdate(false);
            setStatus(() => ({
              isCompleted: true,
              isShowCreatePassWord: false,
              isShowImportWallet: false,
              currentStep: 3,
            }));

            toast.success('✅ Wallet encrypted successfully');
          },
        },
      );
    },
    [walletImport, onStepUpdate],
  );

  const onBackWelcome = useCallback(() => {
    onSwitchModal('Welcome');
  }, [onSwitchModal]);

  const onValidatePrivKey = (e: ChangeEvent<HTMLTextAreaElement>) => {
    const { continueBtn, validationDiv } = refs.current;
    const privKeyEncrypted = e.currentTarget.value;

    if (privKeyEncrypted.length <= 0) {
      if (continueBtn && validationDiv) {
        validationDiv.classList.add('hidden');
        continueBtn.disabled = true;
      }
      return;
    }

    const privKey = DecryptedPrivateKeyFromExport(privKeyEncrypted);

    if (!privKey) {
      if (validationDiv) {
        validationDiv.classList.remove('hidden');
        const validationText = validationDiv.childNodes[0] as HTMLSpanElement;
        if (validationText && continueBtn) {
          continueBtn.disabled = true;
          validationText.textContent = '✗ Invalid private key format';
          validationText.className = 'font-bold text-sm text-red-600';
          return;
        }
      }
      toast.error('Failed. Please try again.');
      return;
    } else {
      if (validationDiv && continueBtn) {
        validationDiv.classList.remove('hidden');
        const validationText = validationDiv.childNodes[0] as HTMLSpanElement;
        if (validationText) {
          continueBtn.disabled = false;
          validationText.textContent = '✓ Valid private key format';
          validationText.className = 'font-bold text-sm text-emerald-600';
          refs.current.privateKey.encode = privKeyEncrypted;
          refs.current.privateKey.value = privKey;
          return;
        }
        toast.error('Failed. Please try again.');
      }
    }
  };

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

  const onBackToHome = async () => {
    await refetch();
    onStepUpdate(false, 1);
    router.push('/');
    router.refresh();
  };

  const onAccessWallet = async () => {
    await refetch();
    onStepUpdate(false, 1);
    router.push('/wallet/me');
    router.refresh();
  };

  useEffect(() => {
    const { continueBtn } = refs.current;

    if (continueBtn) {
      continueBtn.disabled = true;
    }
  }, []);

  return (
    <Fragment>
      <div className="glass-card rounded-3xl p-8">
        <div
          style={{
            animationDelay: '100ms',
          }}
          className="quantum-steps animate-cascase-fade mb-10"
        >
          <div
            className={`quantum-step before:left-[120%] before:h-0.5 ${
              status.isShowImportWallet && 'active'
            } ${status.currentStep > 1 && 'completed'}`}
          >
            <div
              className={`quantum-circle ${
                status.currentStep > 1 && 'completed'
              } w-12 h-12 ${
                status.isShowImportWallet &&
                'active animate-quantum-active-grow'
              }`}
            >
              <span className="text-white font-bold text-sm">1</span>
            </div>
            <p className="text-white text-xs font-semibold mt-4 text-center">
              Private Key
            </p>
          </div>

          <div
            className={`quantum-step before:left-[120%] before:h-0.5 ${
              status.isShowCreatePassWord && 'active'
            } ${status.currentStep > 2 && 'completed'}`}
          >
            <div
              className={`${
                status.currentStep > 2 && 'completed'
              } quantum-circle ${
                status.isShowCreatePassWord &&
                'active animate-quantum-active-grow'
              } w-12 h-12`}
            >
              <span className="text-white font-bold text-sm">2</span>
            </div>
            <p className="text-white text-xs font-semibold mt-4 text-center">
              Create Password
            </p>
          </div>

          <div
            className={`quantum-step  ${
              status.isCompleted && 'active'
            } before:!bg-transparent before:!bg-none`}
          >
            <div
              className={`quantum-circle ${
                status.isCompleted && 'active animate-quantum-active-grow'
              } w-12 h-12`}
            >
              <span className="text-white font-bold text-sm">3</span>
            </div>
            <p className="text-white text-xs font-semibold mt-4 text-center">
              Complete Import
            </p>
          </div>
        </div>

        {status.isShowCreatePassWord && (
          <CreatePassword
            onContinue={onContinueCreatePassword}
            onBack={onBackCreatePassword}
            preview={
              <div className="glass-card rounded-3xl p-4 mb-8">
                <h3 className="text-slate-700 text-center font-bold text-sm mb-4">
                  Wallet Preview
                </h3>
                <div className="key-display-dark rounded-2xl p-6">
                  <div className="text-center text-xs text-slate-400 mb-2">
                    Public Address
                  </div>
                  <code className="text-center block text-xs break-all leading-relaxed">
                    {importWallet.address}
                  </code>
                </div>
              </div>
            }
          />
        )}

        {status.isShowImportWallet && (
          <div className="space-y-4">
            <div
              style={{
                animationDelay: '200ms',
              }}
              className="text-center animate-cascase-fade"
            >
              <h2 className="text-2xl font-black bg-gradient-to-r from-slate-300 via-white to-slate-100 bg-clip-text text-transparent mb-3">
                Import Your Wallet
              </h2>
              <p className="text-slate-200 text-sm font-medium">
                Enter your existing private key to import your wallet
              </p>
            </div>

            <div
              style={{
                animationDelay: '300ms',
              }}
              className="relative enhanced-floating animate-cascase-fade"
            >
              <textarea
                ref={(el) => void (refs.current.textArea = el)}
                className="scrollbar-hide levitating-input w-full rounded-2xl px-6 py-7 text-slate-700 placeholder-transparent font-mono resize-none font-medium text-xs"
                rows={6}
                placeholder=" "
                onInput={onValidatePrivKey}
              ></textarea>
              <label className="text-sm">Private Key (0x...)</label>
            </div>

            <div
              ref={(el) => void (refs.current.validationDiv = el)}
              style={{
                animationDelay: '300ms',
              }}
              className="text-center animate-cascase-fade hidden"
            >
              <span className="font-bold text-sm text-red-600"></span>
            </div>

            <div
              className="flex space-x-6 animate-cascase-fade"
              style={{
                animationDelay: '400ms',
              }}
            >
              <Button onClick={onBackWelcome} variant="glass" size="md">
                ← Back
              </Button>
              <Button
                ref={(el) => void (refs.current.continueBtn = el)}
                variant="quantum"
                size="md"
                onClick={onContinue}
              >
                Continue →
              </Button>
            </div>
          </div>
        )}

        {status.isCompleted && (
          <div>
            <div className="mb-4">
              <div className="w-20 h-20 bg-gradient-to-br from-emerald-500 to-teal-600 rounded-3xl flex items-center justify-center mx-auto mb-4 shadow-2xl">
                <svg
                  className="w-10 h-10 text-white"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path d="M12,1L3,5V11C3,16.55 6.84,21.74 12,23C17.16,21.74 21,16.55 21,11V5L12,1M9,12L7,10L5.5,11.5L9,15L18.5,5.5L17,4L9,12Z"></path>
                </svg>
              </div>

              <h2 className="text-xl text-center font-black text-slate-50 mb-2">
                Wallet Generated!
              </h2>

              <p className="text-slate-100 text-center text-sm leading-relaxed mx-auto">
                Your Web3 wallet is ready. Below are your public credentials —
                keep them safe and share only when necessary.
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
                      {importWallet.address}
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
                    Public Key
                  </label>
                  <span className="text-[10px] text-green-600 bg-green-100 px-2 py-1 rounded-full font-medium">
                    SHAREABLE
                  </span>
                </div>

                <div className="key-display rounded-xl p-4 shimmer">
                  <div className="flex items-center justify-between relative z-10">
                    <code className="text-xs pr-3 text-gray-700 font-mono break-all">
                      {importWallet.pubkey}
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
            </div>

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
        )}
      </div>
    </Fragment>
  );
};

export default ImportWallet;
