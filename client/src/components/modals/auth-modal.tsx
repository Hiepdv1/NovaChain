import { ReAuthModal } from '@/shared/types/modal-type';
import { ModalType, useModalStore } from '@/stores/modal-store';
import { FormEvent, memo, useCallback, useRef, useState } from 'react';
import Input from '../input';
import Button from '../button';
import { DelWalletPool, GetWalletPool } from '@/lib/db/wallet.index';
import { useRouter } from 'next/navigation';
import {
  DecryptPrivateKeyWithPassword,
  SignPayload,
} from '@/lib/crypto/wallet.crypto';
import { buildSignatureWallet } from '@/features/wallet/helpers/wallet-helper';
import { useWalletImport } from '@/features/wallet/hook/useWalletQuery';

interface AuthModalProps extends ReAuthModal {
  closeModal: () => void;
  openModal: (modal: ModalType) => void;
}

const AuthModal = ({
  title,
  des,
  notice,
  wallet,
  closeModal,
}: AuthModalProps) => {
  const inputRef = useRef<HTMLInputElement | null>(null);
  const [isPasswordVisible, setPasswordVisible] = useState(false);
  const [status, setStatus] = useState({
    message: '',
    valid: true,
  });
  const [submit, setSubmit] = useState<{
    isCancel: boolean;
    isConfirm: boolean;
  }>({
    isCancel: false,
    isConfirm: false,
  });
  const router = useRouter();
  const walletImport = useWalletImport();

  const togglePassword = () => {
    const input = inputRef.current;
    setPasswordVisible((prev) => !prev);
    if (input) {
      const start = input.selectionStart ?? input.value.length;
      const end = input.selectionEnd ?? input.value.length;

      requestAnimationFrame(() => {
        if (inputRef.current) {
          const newInput = inputRef.current;
          newInput.focus();
          newInput.setSelectionRange(start, end);
        }
      });
    }
  };

  const onInput = useCallback(() => {
    setStatus((prev) => ({ ...prev, valid: true }));
  }, []);

  const onCancel = useCallback(async () => {
    setSubmit(() => ({
      isCancel: true,
      isConfirm: false,
    }));
    setStatus(() => ({ valid: false, message: '' }));

    await DelWalletPool();

    router.push('/wallet/connect');

    closeModal();

    return;
  }, [router, closeModal]);

  const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setSubmit({
      isCancel: false,
      isConfirm: true,
    });
    const form = e.currentTarget;
    const formData = new FormData(form);
    const data = Object.fromEntries(formData) as { password: string };

    const walletPool = await GetWalletPool();

    if (walletPool.length < 1) {
      router.push('/wallet/connect');
      return;
    }
    const wallet = walletPool[0];

    const privateKey = DecryptPrivateKeyWithPassword(
      data.password,
      wallet.encryptedPrivateKey.cipherText,
      wallet.encryptedPrivateKey.iv,
      wallet.encryptedPrivateKey.authTag,
    );

    if (!privateKey) {
      setStatus({
        valid: false,
        message: 'Incorrect password. Please try again.',
      });
      setSubmit({
        isCancel: false,
        isConfirm: false,
      });
      return;
    } else {
      setStatus({ valid: true, message: '' });
    }

    const payload = buildSignatureWallet(privateKey);

    const sig = SignPayload(privateKey, payload);

    walletImport.mutate(
      {
        data: payload,
        sig,
      },
      {
        onError: (err) => {
          setSubmit({
            isCancel: false,
            isConfirm: false,
          });
          setStatus({
            message: err.message,
            valid: false,
          });
        },
        onSuccess: () => {
          setSubmit({
            isCancel: false,
            isConfirm: false,
          });
          useModalStore.getState().actions.closeModal();
          router.refresh();
        },
      },
    );
  };

  return (
    <div className="z-[9999] fixed inset-0 bg-black/40 flex items-center justify-center p-4">
      <div className="overflow-hidden bg-[rgba(255,255,255,.9)] backdrop-blur-xl max-w-md mx-auto w-full rounded-3xl">
        <div className="p-8">
          <div className="text-center mb-8">
            <div className="w-16 h-16 bg-primary shadow-[0_8px_25px_rgba(251,113,133,0.3)] rounded-2xl flex items-center justify-center mx-auto mb-6">
              <svg
                className="w-12 h-12 text-white"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="1.5"
                  d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"
                ></path>
              </svg>
            </div>

            <h2 className="text-xl font-semibold text-gray-800 mb-3">
              {title}
            </h2>

            <p className="text-gray-600 text-[13px] leading-relaxed">{des}</p>
          </div>

          <div className="mb-8">
            <div className="glass-card rounded-2xl p-6 border border-white/20 relative overflow-hidden">
              <div className="absolute inset-0 bg-[lab(100_0_0)]" />
              <div className="relative z-10 flex items-center space-x-4">
                <div className="shrink-0">
                  <div className="animate-glass-float rounded-xl flex items-center justify-center shadow-xl w-12 h-12 bg-primary">
                    <svg
                      className="w-7 h-7 text-white"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z"
                      ></path>
                    </svg>
                  </div>
                </div>

                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between mb-3">
                    <h3 className="text-xs font-bold text-black/90 uppercase tracking-wider">
                      Premium Wallet
                    </h3>
                    <div className="flex items-center space-x-2">
                      <div className="animate-pulse-fade w-2 h-2 bg-green-500 rounded-full" />
                      <div className="text-[10px] text-green-500 font-semibold">
                        SECURE
                      </div>
                    </div>
                  </div>

                  <div className="glass-card rounded-lg p-3 border border-white/10">
                    <p className="text-xs break-all text-black">
                      {wallet.address}
                    </p>
                  </div>

                  <div className="flex items-center justify-end mt-3">
                    <span className="text-xs font-bold text-black/70">
                      Network:
                      <span className="text-xs text-primary"> CCC</span>
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <form onSubmit={onSubmit} className="space-y-6">
            <div>
              {/* <label
                htmlFor="password"
                className="block text-sm font-medium text-gray-400 mb-2"
              >
                Password
              </label> */}
              <div className="relative">
                <div className="enhanced-floating">
                  <Input
                    onInput={onInput}
                    id="password"
                    name="password"
                    type={isPasswordVisible ? 'text' : 'password'}
                    placeholder=" "
                    variant="levitating"
                    inputSize="md"
                    className={`${
                      !status.valid ? '!border-red-600' : ''
                    } text-xs peer`}
                    ref={(el) => void (inputRef.current = el)}
                    disabled={submit.isConfirm || submit.isCancel}
                  />
                  <label htmlFor="password" className="text-sm">
                    Password
                  </label>
                </div>
                <button
                  type="button"
                  className="cursor-pointer absolute right-4 top-1/2 transform -translate-y-1/2 text-gray-400 transition-colors duration-200"
                  tabIndex={-1}
                  aria-label="Toggle password visibility"
                  onClick={togglePassword}
                >
                  {isPasswordVisible ? (
                    <svg
                      id="eyeIcon"
                      className="w-5 h-5"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21"
                      ></path>
                    </svg>
                  ) : (
                    <svg
                      id="eyeIcon"
                      className="w-5 h-5 peer-hover:text-primary"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"
                      ></path>
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                      ></path>
                    </svg>
                  )}
                </button>
              </div>

              {submit.isConfirm && (
                <div className="mt-3">
                  <p className="text-xs text-gray-600 font-semibold mt-2 text-center">
                    Verifying your credentials...
                  </p>
                </div>
              )}

              {!status.valid && status.message && (
                <div className="mt-3 p-3 bg-red-50 border-red-200 border rounded-xl">
                  <div className="flex items-center text-red-600">
                    <svg
                      className="w-4 h-4 mr-2"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
                        clipRule="evenodd"
                      ></path>
                    </svg>
                    <span className="text-xs font-medium">
                      {status.message}
                    </span>
                  </div>
                </div>
              )}
            </div>

            <div className="flex space-x-4">
              <Button
                type="button"
                className="py-4 text-sm disabled:opacity-50 flex justify-center items-center"
                variant="default"
                size="md"
                onClick={onCancel}
                disabled={!status.valid}
              >
                Cancel
                {submit.isCancel && (
                  <svg
                    id="loadingIcon"
                    className="animate-spin w-4 h-4 ml-2"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    ></circle>
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                )}
              </Button>
              <Button
                type="submit"
                className="flex items-center justify-center py-4 text-sm disabled:opacity-50 flex-1"
                variant="secondary"
                size="md"
                disabled={!status.valid}
              >
                <span>Confirm</span>
                {submit.isConfirm && (
                  <svg
                    id="loadingIcon"
                    className="animate-spin w-4 h-4 ml-2"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    ></circle>
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                )}
              </Button>
            </div>
          </form>

          {notice}
        </div>
      </div>
    </div>
  );
};

export default memo(AuthModal);
