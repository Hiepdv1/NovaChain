/* eslint-disable @typescript-eslint/no-explicit-any */
'use client';

import Button from '@/components/button';
import Input from '@/components/input';
import useWalletContext from '@/components/providers/wallet-provider';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { ValidateAddress } from '@/lib/db/wallet.store';
import { Loader } from 'lucide-react';

import { useRouter } from 'next/navigation';
import {
  ChangeEvent,
  FormEvent,
  useCallback,
  useEffect,
  useRef,
  useState,
} from 'react';
import TransactionFee, {
  FeeList,
  FeeOption,
} from '../components/transaction-fee';
import TransactionMessage from '../components/transaction-message';
import { IsValidNumber } from '@/lib/utils';
import ContentLoading from '@/components/loading/content-loading';
import { toast } from '@/components/globalToaster';
import TransactionPreview from '../components/transaction-preview';
import NetworkStatus from '../components/network-status';
import { useModalStore } from '@/stores/modal-store';
import { StoredWallet } from '@/shared/types/wallet';
import { CreateNewTXPayload, Transaction } from '../types/transaction';

type RefMap = {
  btnVerify: HTMLButtonElement | null;
  inputRecipient: HTMLInputElement | null;
  reviewBtn: HTMLButtonElement | null;
  amountInput: HTMLInputElement | null;
};

type ValidateForm = {
  recipientAddr: {
    valid: boolean;
    isLoading: boolean;
    message: string | null;
    isLock: boolean;
    value: string;
  };
  amount: {
    valid: boolean;
    message: string | null;
    value: number;
  };
  message?: string;
};

type SubmitForm = {
  amount: string;
  fee: string;
  message: string;
  from: string;
  to: string;
};

const SendTransactionPage = () => {
  const router = useRouter();
  const { wallet } = useWalletContext();
  const { actions } = useModalStore();
  const refs = useRef<RefMap>({
    btnVerify: null,
    inputRecipient: null,
    reviewBtn: null,
    amountInput: null,
  });
  const [validate, setValidate] = useState<ValidateForm>({
    recipientAddr: {
      isLoading: false,
      valid: false,
      isLock: false,
      message: null,
      value: '',
    },
    amount: {
      message: null,
      valid: false,
      value: 0,
    },
    message: '',
  });
  const [currentFee, setCurrentFee] = useState<FeeOption>(
    FeeList.find((f) => f.checked)!,
  );

  const onCheckedFee = useCallback((e: ChangeEvent<HTMLInputElement>) => {
    const input = e.currentTarget;
    const value = input.value;

    const fee = FeeList.find((f) => f.value === value);
    if (!fee) return;

    setCurrentFee(fee);
  }, []);

  const onVerifyAddress = async () => {
    if (!wallet || !wallet.data) return;
    setValidate((prev) => {
      return {
        ...prev,
        recipientAddr: {
          isLoading: true,
          valid: false,
          message: null,
          isLock: false,
          value: '',
        },
      };
    });
    const { inputRecipient, btnVerify } = refs.current;

    if (!inputRecipient || !btnVerify) return;
    btnVerify.disabled = true;
    const value = inputRecipient.value;

    if (ValidateAddress(value)) {
      if (value.trim() === wallet.data.Address.trim()) {
        setValidate((prev) => {
          return {
            ...prev,
            recipientAddr: {
              isLoading: false,
              valid: false,
              isLock: false,
              message: `The recipient address cannot be the same as your wallet address.`,
              value: '',
            },
          };
        });
        return;
      }
      setValidate((prev) => {
        return {
          ...prev,
          recipientAddr: {
            isLoading: false,
            valid: true,
            isLock: true,
            message: `✓ Address verified and locked • Active account • Click "Edit" to change`,
            value: value.trim(),
          },
        };
      });
    } else {
      setValidate((prev) => {
        return {
          ...prev,
          recipientAddr: {
            isLoading: false,
            valid: false,
            isLock: false,
            message: `❌ Invalid address format. Must be a valid Base58 string (letters A–Z, a–z, and numbers 1–9, excluding 0, O, I, l).`,
            value: '',
          },
        };
      });
    }
  };

  const onValidateAddr = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { btnVerify } = refs.current;
    if (!btnVerify) return;

    setValidate((prev) => {
      return {
        ...prev,
        recipientAddr: {
          isLoading: false,
          isLock: false,
          valid: false,
          message: null,
          value: '',
        },
      };
    });

    const input = e.currentTarget;
    const value = input.value;

    if (value.length === 34) {
      btnVerify.disabled = false;
    } else {
      btnVerify.disabled = true;
    }
  };

  const onEditAddress = () => {
    setValidate((prev) => {
      return {
        ...prev,
        recipientAddr: {
          isLoading: false,
          isLock: false,
          valid: false,
          message: null,
          value: '',
        },
      };
    });
  };

  const onInputAmount = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (!wallet?.data) return;

    const input = e.currentTarget;
    const value = input.value;

    if (value === '') {
      setValidate((prev) => ({
        ...prev,
        amount: {
          message: null,
          valid: false,
          value: 0,
        },
      }));
      return;
    }

    if (IsValidNumber(value)) {
      const amount = parseFloat(value);

      if (amount <= 0) {
        setValidate((prev) => ({
          ...prev,
          amount: {
            message: 'Amount must be greater than 0',
            valid: false,
            value: 0,
          },
        }));
      } else if (amount > parseFloat(wallet.data.Balance)) {
        setValidate((prev) => ({
          ...prev,
          amount: {
            message: `Insufficient balance. Maximum available: ${wallet.data.Balance} CCC`,
            valid: false,
            value: 0,
          },
        }));
      } else {
        setValidate((prev) => ({
          ...prev,
          amount: {
            message: null,
            valid: true,
            value: amount,
          },
        }));
      }
    } else {
      setValidate((prev) => ({
        ...prev,
        amount: {
          message: 'Invalid number format',
          valid: false,
          value: 0,
        },
      }));
    }
  };

  const onConfirmPassword = (
    w: StoredWallet,
    privKey: string,
    data: Transaction,
    ...args: any
  ) => {
    const payload = args[0] as CreateNewTXPayload;
    actions.openModal({
      type: 'previewTx',
      props: {
        amount: payload.data.amount,
        fee: payload.data.fee,
        message: payload.data.message,
        to: payload.data.to,
        from: w.address,
        balance: parseFloat(wallet?.data?.Balance || '0'),
        transactions: data,
      },
    });
  };

  const onSubmit = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    if (!wallet || !wallet.data) return;

    const { inputRecipient } = refs.current;

    if (!inputRecipient) return;

    const formData = new FormData(e.currentTarget);
    const data = Object.fromEntries(formData) as SubmitForm;

    const fee = FeeList.find((f) => f.value === data.fee);

    if (!fee) return;

    if (fee.fee + Number(data.amount) > Number(wallet.data.Balance)) {
      toast.warning('You dont have enough amount.');
      return;
    }

    actions.openModal({
      type: 'verifyTx',
      props: {
        onSubmit: onConfirmPassword,
        data: {
          amount: Number(data.amount),
          fee: fee.fee,
          timestamp: Math.floor(Date.now() / 1000),
          to: inputRecipient.value,
          message: data.message,
        },
      },
    });
  };

  const onClickSetAmount = (percentage: number) => {
    const { amountInput } = refs.current;
    if (!amountInput || !wallet?.data) return;

    amountInput.value = (
      parseFloat(wallet.data.Balance) *
      (percentage / 100)
    ).toString();

    const event = new Event('input', { bubbles: true });

    amountInput.dispatchEvent(event);
    toast.info(
      'Amount Set',
      `Set to ${percentage}% of balance (${wallet.data.Balance} CCC)`,
    );
  };

  const onInputMessage = useCallback((e: ChangeEvent<HTMLTextAreaElement>) => {
    const textarea = e.currentTarget;
    setValidate((prev) => {
      return {
        ...prev,
        message: textarea.value,
      };
    });
  }, []);

  useEffect(() => {
    const updateFormButton = () => {
      const { reviewBtn } = refs.current;
      if (!reviewBtn) return;

      if (validate.recipientAddr.valid && validate.amount.valid) {
        reviewBtn.disabled = false;
      } else {
        reviewBtn.disabled = true;
      }
    };

    updateFormButton();
  }, [validate]);

  useEffect(() => {
    const { btnVerify, reviewBtn } = refs.current;

    if (btnVerify) {
      btnVerify.disabled = true;
    }

    if (reviewBtn) {
      reviewBtn.disabled = true;
    }
  }, [refs]);

  useEffect(() => {
    if (!wallet || wallet.error || !wallet.data) {
      router.push('/');
    }
  }, [wallet, router]);

  if (!wallet || !wallet.data) {
    return <ContentLoading />;
  }

  return (
    <div className="space-y-6 select-none">
      <div className="space-y-6">
        <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between">
          <div>
            <h2 className="text-xl font-bold text-gray-900 dark:text-white mb-2">
              Send Transaction
            </h2>
            <p className="text-xs text-gray-700 dark:text-gray-400">
              Send CCC tokens securely on the CryptoChain network
            </p>
          </div>

          <div className="flex flex-col lg:flex-row items-center space-x-4">
            <div className="text-center">
              <div className="text-sm text-gray-500 dark:text-gray-400">
                Block Time
              </div>
              <div className="text-sm font-semibold text-gray-900 dark:text-white">
                12.3s
              </div>
            </div>

            <div className="text-center">
              <div className="text-sm text-gray-500 dark:text-gray-400">
                Network Load
              </div>
              <div className="text-sm font-semibold text-green-600 dark:text-green-400">
                Normal
              </div>
            </div>
          </div>
        </div>

        <div className="glass-card dark:border-secondary-dark dark:bg-primary-dark rounded-2xl p-4 border-l-4 border-green-500">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse" />
              <div>
                <div className="text-sm font-medium text-gray-900 dark:text-white">
                  Network Status: Operational
                </div>
                <div className="text-xs text-gray-600 dark:text-gray-400">
                  All systems running normally • 1,847 pending transactions
                </div>
              </div>
            </div>

            <div className="text-right">
              <div className="text-xs text-gray-500 dark:text-gray-400">
                Last Block
              </div>
              <div className="text-sm font-medium text-gray-900 dark:text-white">
                #2,847,392
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 xl:grid-cols-4 gap-8">
        <div className="xl:col-span-2">
          <div className="glass-card dark:bg-primary-dark dark:border-secondary-dark rounded-2xl p-6 shadow-lg shadow-black/5 dark:shadow-black/20">
            <div className="flex items-center justify-between mb-6">
              <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                Transaction Details
              </h2>
              <div className="flex items-center space-x-2 text-xs text-gray-500 dark:text-gray-400">
                <svg
                  className="w-3 h-3"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path
                    fillRule="evenodd"
                    d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
                    clipRule="evenodd"
                  ></path>
                </svg>

                <span>Secure Transaction</span>
              </div>
            </div>

            <form onSubmit={onSubmit} className="space-y-6 select-none">
              <div className="space-y-2">
                <label
                  className="text-xs font-bold cursor-text block"
                  htmlFor="fromAddr"
                >
                  From Address (Your wallet)
                </label>

                <div className="relative w-full">
                  <Input
                    className="text-gray-500 cursor-not-allowed text-xs py-3.5 dark:bg-primary-dark focus:!bg-transparent"
                    variant="levitating"
                    inputSize="sm"
                    name="fromAddr"
                    id="fromAddr"
                    disabled
                    value={wallet?.data?.Address}
                  />

                  <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
                    <svg
                      className="w-5 h-5 text-green-500"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                        clipRule="evenodd"
                      ></path>
                    </svg>
                  </div>
                </div>

                <div className="flex items-center text-xs">
                  <span className="text-gray-500 dark:text-gray-400">
                    Balance:
                    <span className="ml-1 font-medium text-gray-900 dark:text-white">
                      {wallet?.data?.Balance} CCC
                    </span>
                  </span>
                </div>
              </div>

              <div className="space-y-2">
                <div className="flex items-center space-x-2">
                  <label
                    className="text-xs font-bold cursor-text block"
                    htmlFor="recipientAddr"
                  >
                    Recipient Address
                  </label>

                  <Tooltip>
                    <TooltipTrigger asChild>
                      <svg
                        className="w-4 h-4 text-gray-400 cursor-help"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-3a1 1 0 00-.867.5 1 1 0 11-1.731-1A3 3 0 0113 8a3.001 3.001 0 01-2 2.83V11a1 1 0 11-2 0v-1a1 1 0 011-1 1 1 0 100-2zm0 8a1 1 0 100-2 1 1 0 000 2z"
                          clipRule="evenodd"
                        ></path>
                      </svg>
                    </TooltipTrigger>
                    <TooltipContent
                      align="center"
                      className="bg-gray-600 dark:bg-white/95 p-4"
                    >
                      <strong>Address Validation</strong>
                      <ul className="list-disc list-inside mt-2 space-y-1 text-gray-300 dark:text-gray-950">
                        <li>
                          Must be a valid CryptoChain address (40 characters)
                        </li>
                        <li>Format and length are checked automatically</li>
                        <li>
                          Click <strong>Verify</strong> to ensure the address
                          exists on the network
                        </li>
                        <li>Transactions are irreversible once submitted</li>
                        <li>Always double-check before sending</li>
                      </ul>
                    </TooltipContent>
                  </Tooltip>
                </div>

                <div className="flex space-x-3">
                  <div className="relative flex-1 ">
                    <Input
                      className={`${
                        !validate.recipientAddr.valid &&
                        validate.recipientAddr.message
                          ? '!border-[#ef4444] !shadow-none'
                          : ''
                      } ${
                        validate.recipientAddr.valid
                          ? '!border-green-600 !shadow-none'
                          : ''
                      } text-gray-900 !transform-none placeholder:text-gray-700 dark:text-white dark:placeholder:text-gray-400 text-xs py-3.5 dark:bg-primary-dark focus:!bg-transparent pr-12`}
                      variant="levitating"
                      inputSize="sm"
                      name="recipientAddr"
                      id="recipientAddr"
                      onInput={onValidateAddr}
                      disabled={validate.recipientAddr.isLock}
                      ref={(el) => void (refs.current.inputRecipient = el)}
                      placeholder="Enter recipient address (34 characters)"
                    />
                    <div className="absolute right-3 top-1/2 transform -translate-y-1/2">
                      {validate.recipientAddr.message && (
                        <>
                          {!validate.recipientAddr.valid && (
                            <svg
                              id="addressInvalid"
                              className="w-5 h-5 text-red-500"
                              fill="currentColor"
                              viewBox="0 0 20 20"
                            >
                              <path
                                fillRule="evenodd"
                                d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                                clipRule="evenodd"
                              ></path>
                            </svg>
                          )}

                          {validate.recipientAddr.valid && (
                            <svg
                              className="w-5 h-5 text-green-600 dark:text-green-400"
                              fill="currentColor"
                              viewBox="0 0 20 20"
                            >
                              <path
                                fillRule="evenodd"
                                d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
                                clipRule="evenodd"
                              ></path>
                            </svg>
                          )}
                        </>
                      )}

                      {validate.recipientAddr.isLoading && (
                        <Loader className="w-5 h-5 animate-spin" />
                      )}
                    </div>
                  </div>

                  {validate.recipientAddr.isLock ? (
                    <Button
                      variant="default"
                      type="button"
                      size="sm"
                      disabled={false}
                      onClick={onEditAddress}
                      className="!bg-gray-700 hover:!bg-gray-600 !text-white !border-primary-dark flex space-x-1 items-center justify-center flex-none w-auto text-xs shadow-none animate-none !transform-none py-2 px-6"
                    >
                      <svg
                        className="w-4 h-4"
                        fill="none"
                        stroke="currentColor"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth="2"
                          d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                        ></path>
                      </svg>
                      <span>Edit</span>
                    </Button>
                  ) : (
                    <Button
                      ref={(el) => void (refs.current.btnVerify = el)}
                      type="button"
                      className="text-xs w-auto shadow-none animate-none transform-none py-2 px-6"
                      size="sm"
                      variant="secondary"
                      onClick={onVerifyAddress}
                    >
                      {validate.recipientAddr.isLoading
                        ? 'Waiting...'
                        : 'Verify'}
                    </Button>
                  )}
                </div>

                {validate.recipientAddr.message && (
                  <div
                    className={`text-xs ${
                      validate.recipientAddr.valid
                        ? 'text-green-600 dark:text-green-400'
                        : 'text-red-600 dark:text-red-400'
                    } flex items-center space-x-2`}
                  >
                    {validate.recipientAddr.valid ? (
                      <svg
                        className="w-4 h-4"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                          clipRule="evenodd"
                        ></path>
                      </svg>
                    ) : (
                      <svg
                        className="w-4 h-4"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                          clipRule="evenodd"
                        ></path>
                      </svg>
                    )}

                    <span>{validate.recipientAddr.message}</span>
                  </div>
                )}
              </div>

              <div className="space-y-2">
                <div className="flex items-center space-x-2">
                  <label
                    className="block text-xs font-medium text-gray-900 dark:text-white"
                    htmlFor="amount"
                  >
                    Amount to Send
                  </label>
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <svg
                        className="w-4 h-4 text-gray-400 cursor-help"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-3a1 1 0 00-.867.5 1 1 0 11-1.731-1A3 3 0 0113 8a3.001 3.001 0 01-2 2.83V11a1 1 0 11-2 0v-1a1 1 0 011-1 1 1 0 100-2zm0 8a1 1 0 100-2 1 1 0 000 2z"
                          clipRule="evenodd"
                        ></path>
                      </svg>
                    </TooltipTrigger>
                    <TooltipContent
                      align="center"
                      className="bg-gray-600 dark:bg-white/95 p-4"
                    >
                      <strong>Amount Guidelines</strong>

                      <ul className="list-disc list-inside mt-2 space-y-1 text-gray-300 dark:text-gray-950">
                        <li>Minimum: 0.00000001 CCC</li>
                        <li>Maximum: Your available balance minus fees</li>
                        <li>
                          Precision: <strong>Up to 8 decimal places</strong>
                        </li>
                        <li>
                          Always keep some CCC for future transaction fees
                        </li>
                      </ul>
                    </TooltipContent>
                  </Tooltip>
                </div>

                <div className="relative">
                  <Input
                    className={`${
                      !validate.amount.valid && validate.amount.message
                        ? '!border-[#ef4444] !shadow-none'
                        : ''
                    } ${
                      validate.amount.valid
                        ? '!border-green-600 !shadow-none'
                        : ''
                    } text-gray-900 !transform-none placeholder:text-gray-700 dark:text-white dark:placeholder:text-gray-400 text-xs py-3.5 dark:bg-primary-dark focus:!bg-transparent pr-14`}
                    variant="levitating"
                    inputSize="sm"
                    name="amount"
                    id="amount"
                    type="text"
                    onInput={onInputAmount}
                    min={0}
                    max={wallet?.data?.Balance}
                    ref={(el) => void (refs.current.amountInput = el)}
                    placeholder="0.00000000"
                  />

                  <div className="absolute right-2 top-1/2 transform -translate-1/2 text-sm font-medium text-gray-500 dark:text-gray-400">
                    CCC
                  </div>
                </div>

                {!validate.amount.valid && validate.amount.message && (
                  <div
                    className={`text-xs ${
                      !validate.amount.valid
                        ? 'text-red-600 dark:text-red-400'
                        : ''
                    } flex items-center space-x-2`}
                  >
                    <svg
                      className="w-4 h-4"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                        clipRule="evenodd"
                      ></path>
                    </svg>
                    <span>{validate.amount.message}</span>
                  </div>
                )}

                <div className="mt-4 grid grid-cols-4 gap-2">
                  {[25, 50, 75, 100].map((value) => {
                    return (
                      <Button
                        type="button"
                        key={value}
                        onClick={() => onClickSetAmount(value)}
                        className="font-bold bg-primary-50 dark:bg-primary-900/20 hover:bg-primary-100 dark:hover:bg-primary-900/30 !translate-none text-xs text-primary-600 dark:text-primary-400 !transition-colors border-none "
                        variant="default"
                        size="sm"
                      >
                        {value}%
                      </Button>
                    );
                  })}
                </div>
              </div>

              {/* Transaction Fee */}
              <TransactionFee onChecked={onCheckedFee} />

              {/* Transaction Message */}
              <TransactionMessage onInputMessage={onInputMessage} />

              {/* Transaction actions */}
              <div className="flex flex-col sm:flex-row space-y-3 sm:space-y-0 sm:space-x-4 pt-6">
                <Button
                  className="flex items-center justify-center space-x-2"
                  size="md"
                  variant="secondary"
                  type="submit"
                  ref={(el) => void (refs.current.reviewBtn = el)}
                >
                  <svg
                    className="w-5 h-5"
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
                  <span>Review & Send Transaction</span>
                </Button>
              </div>
            </form>
          </div>
        </div>

        <div className="xl:col-span-2 space-y-6">
          <TransactionPreview
            recipient={
              validate.recipientAddr.valid
                ? refs.current.inputRecipient?.value
                : ''
            }
            active={validate.recipientAddr.valid && validate.amount.valid}
            fromAddr={wallet.data.Address.trim()}
            data={{
              amount: validate.amount.value,
              fee: currentFee.fee,
              message: validate.message || '',
            }}
          />

          <NetworkStatus />
        </div>
      </div>
    </div>
  );
};

export default SendTransactionPage;
