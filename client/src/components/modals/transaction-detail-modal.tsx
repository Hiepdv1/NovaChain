'use client';

import { formatAddress } from '@/lib/utils';
import { TransactionDetailModalProps } from '@/shared/types/modal-type';
import { useModalStore } from '@/stores/modal-store';
import { Loader2 } from 'lucide-react';
import { useState } from 'react';

const TransactionDetailModal = ({
  balance,
  from,
  to,
  amount,
  message,
  fee,
}: TransactionDetailModalProps) => {
  const { actions } = useModalStore();
  const [isLoading, setIsloading] = useState(false);

  const onSendTransaction = () => {
    setIsloading(true);
  };

  return (
    <div className="fixed z-[999] inset-0 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4">
      <div className="glass-card bg-white dark:border-secondary-dark dark:bg-primary-dark rounded-2xl p-4 sm:p-6 lg:p-8 w-full max-w-sm sm:max-w-2xl lg:max-w-4xl shadow-2xl max-h-[90vh] sm:max-h-[95vh] overflow-y-auto">
        <div className="text-center mb-4 sm:mb-6 lg:mb-8">
          <div className="w-12 h-12 sm:w-16 sm:h-16 lg:w-20 lg:h-20 bg-gradient-primary rounded-full flex items-center justify-center mx-auto mb-3 sm:mb-4">
            <svg
              className="w-6 h-6 sm:w-8 sm:h-8 lg:w-10 lg:h-10 text-white"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              ></path>
            </svg>
          </div>
          <h3 className="text-xl sm:text-2xl lg:text-3xl font-bold text-gray-900 dark:text-white mb-2 sm:mb-3">
            Review Transaction Details
          </h3>
          <p className="text-sm sm:text-base text-gray-600 dark:text-gray-400 px-2">
            Please carefully review all details before confirming this
            transaction
          </p>
          <div className="mt-3 sm:mt-4 inline-flex items-center space-x-2 px-3 sm:px-4 py-1.5 sm:py-2 bg-blue-50 dark:bg-blue-900/20 rounded-full">
            <svg
              className="w-3 h-3 sm:w-4 sm:h-4 text-blue-600 dark:text-blue-400"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
                clipRule="evenodd"
              ></path>
            </svg>
            <span className="text-xs sm:text-sm font-medium text-blue-800 dark:text-blue-200">
              Secure &amp; Encrypted
            </span>
          </div>
        </div>

        <div className="space-y-4 sm:space-y-6 lg:space-y-8">
          <div className="p-4 sm:p-6 lg:p-8 bg-gradient-to-r from-blue-50 to-purple-50 dark:from-blue-900/20 dark:to-purple-900/20 rounded-xl sm:rounded-2xl border border-blue-200/50 dark:border-blue-700/50">
            <h4 className="text-base sm:text-lg font-semibold text-gray-900 dark:text-white mb-4 sm:mb-6 text-center">
              Transaction Flow
            </h4>

            <div className="block sm:hidden space-y-4">
              <div className="text-center">
                <div className="w-16 h-16 bg-gradient-primary rounded-full flex items-center justify-center mb-3 mx-auto shadow-lg">
                  <svg
                    className="w-8 h-8 text-white"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M4 4a2 2 0 00-2 2v1h16V6a2 2 0 00-2-2H4zM18 9H2v5a2 2 0 002 2h12a2 2 0 002-2V9zM4 13a1 1 0 011-1h1a1 1 0 110 2H5a1 1 0 01-1-1zm5-1a1 1 0 100 2h1a1 1 0 100-2H9z"></path>
                  </svg>
                </div>
                <div className="text-sm font-bold text-gray-900 dark:text-white mb-2">
                  From (Your Wallet)
                </div>
                <div className="text-xs text-gray-500 dark:text-gray-400 font-mono bg-white/50 dark:bg-white/10 px-2 py-1 rounded-lg mb-2 break-all">
                  {formatAddress(from)}
                </div>
                <div className="text-xs text-gray-600 dark:text-gray-300">
                  Balance: {balance} CCC
                </div>
              </div>

              <div className="flex justify-center">
                <div className="text-center bg-white dark:bg-gray-800 p-3 rounded-xl shadow-lg border border-gray-200 dark:border-gray-600">
                  <svg
                    className="w-6 h-6 text-primary-600 dark:text-primary-400 mx-auto mb-1"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M19 14l-7 7m0 0l-7-7m7 7V3"
                    ></path>
                  </svg>
                  <div className="text-sm font-bold text-gray-900 dark:text-white">
                    {amount} CCC
                  </div>
                  <div className="text-xs text-gray-500 dark:text-gray-400">
                    Sending
                  </div>
                </div>
              </div>

              <div className="text-center">
                <div className="w-16 h-16 bg-gray-200 dark:bg-gray-700 rounded-full flex items-center justify-center mb-3 mx-auto shadow-lg">
                  <svg
                    className="w-8 h-8 text-gray-500 dark:text-gray-400"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z"></path>
                  </svg>
                </div>
                <div className="text-sm font-bold text-gray-900 dark:text-white mb-2">
                  To (Recipient)
                </div>
                <div className="text-xs text-gray-500 dark:text-gray-400 font-mono bg-white/50 dark:bg-white/10 px-2 py-1 rounded-lg mb-2 break-all">
                  {formatAddress(to)}
                </div>
                <div className="flex items-center justify-center space-x-1">
                  <svg
                    className="w-3 h-3 text-green-500"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                      clipRule="evenodd"
                    ></path>
                  </svg>
                  <span className="text-xs text-green-600 dark:text-green-400 font-medium">
                    Verified
                  </span>
                </div>
              </div>
            </div>

            <div className="hidden sm:flex items-center justify-between">
              <div className="text-center flex-1">
                <div className="w-16 h-16 lg:w-20 lg:h-20 bg-gradient-primary rounded-full flex items-center justify-center mb-3 lg:mb-4 mx-auto shadow-lg">
                  <svg
                    className="w-8 h-8 lg:w-10 lg:h-10 text-white"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M4 4a2 2 0 00-2 2v1h16V6a2 2 0 00-2-2H4zM18 9H2v5a2 2 0 002 2h12a2 2 0 002-2V9zM4 13a1 1 0 011-1h1a1 1 0 110 2H5a1 1 0 01-1-1zm5-1a1 1 0 100 2h1a1 1 0 100-2H9z"></path>
                  </svg>
                </div>
                <div className="text-sm lg:text-lg font-bold text-gray-900 dark:text-white mb-2">
                  From (Your Wallet)
                </div>
                <div className="text-xs lg:text-sm text-gray-500 dark:text-gray-400 font-mono bg-white/50 dark:bg-white/10 px-2 lg:px-3 py-1 lg:py-2 rounded-lg mb-2 break-all">
                  {from}
                </div>
                <div className="text-xs lg:text-sm text-gray-600 dark:text-gray-300">
                  Balance: {balance} CCC
                </div>
              </div>

              <div className="flex-1 flex items-center justify-center px-4 lg:px-8">
                <div className="flex items-center space-x-2 lg:space-x-4">
                  <div className="w-8 lg:w-16 h-0.5 lg:h-1 bg-gradient-primary rounded-full"></div>
                  <div className="text-center bg-white dark:bg-gray-800 p-2 lg:p-4 rounded-xl lg:rounded-2xl shadow-lg border border-gray-200 dark:border-gray-600">
                    <svg
                      className="w-6 h-6 lg:w-8 lg:h-8 text-primary-600 dark:text-primary-400 mx-auto mb-1 lg:mb-2"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth="2"
                        d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"
                      ></path>
                    </svg>
                    <div className="text-sm lg:text-xl font-bold text-gray-900 dark:text-white">
                      {amount} CCC
                    </div>
                    <div className="text-xs text-gray-500 dark:text-gray-400">
                      Sending
                    </div>
                  </div>
                  <div className="w-8 lg:w-16 h-0.5 lg:h-1 bg-gradient-primary rounded-full"></div>
                </div>
              </div>

              <div className="text-center flex-1">
                <div className="w-16 h-16 lg:w-20 lg:h-20 bg-gray-200 dark:bg-gray-700 rounded-full flex items-center justify-center mb-3 lg:mb-4 mx-auto shadow-lg">
                  <svg
                    className="w-8 h-8 lg:w-10 lg:h-10 text-gray-500 dark:text-gray-400"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z"></path>
                  </svg>
                </div>
                <div className="text-sm lg:text-lg font-bold text-gray-900 dark:text-white mb-2">
                  To (Recipient)
                </div>
                <div className="text-xs lg:text-sm text-gray-500 dark:text-gray-400 font-mono bg-white/50 dark:bg-white/10 px-2 lg:px-3 py-1 lg:py-2 rounded-lg mb-2 break-all">
                  {to}
                </div>
                <div className="text-xs text-gray-500 dark:text-gray-400">
                  {formatAddress(to)}
                </div>
                <div className="flex items-center justify-center space-x-1 mt-2">
                  <svg
                    className="w-3 h-3 lg:w-4 lg:h-4 text-green-500"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                      clipRule="evenodd"
                    ></path>
                  </svg>
                  <span className="text-xs font-medium text-green-600 dark:text-green-400">
                    Address Verified
                  </span>
                </div>
              </div>
            </div>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 sm:gap-6">
            <div className="lg:col-span-2 space-y-4 sm:space-y-6">
              <div className="bg-white/70 dark:bg-gray-800/70 rounded-xl sm:rounded-2xl p-4 sm:p-6 border border-gray-200/50 dark:border-gray-600/50">
                <h4 className="text-lg sm:text-xl font-semibold text-gray-900 dark:text-white mb-4 sm:mb-6 flex items-center space-x-2">
                  <svg
                    className="w-5 h-5 sm:w-6 sm:h-6 text-primary-600 dark:text-primary-400"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M4 4a2 2 0 00-2 2v1h16V6a2 2 0 00-2-2H4zM18 9H2v5a2 2 0 002 2h12a2 2 0 002-2V9zM4 13a1 1 0 011-1h1a1 1 0 110 2H5a1 1 0 01-1-1zm5-1a1 1 0 100 2h1a1 1 0 100-2H9z"></path>
                  </svg>
                  <span>Financial Summary</span>
                </h4>

                <div className="space-y-3 sm:space-y-4">
                  <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center p-3 sm:p-5 bg-blue-50 dark:bg-blue-900/20 rounded-xl border border-blue-200/50 dark:border-blue-700/50 space-y-3 sm:space-y-0">
                    <div className="flex items-center space-x-3">
                      <div className="w-8 h-8 sm:w-10 sm:h-10 bg-blue-500 rounded-full flex items-center justify-center flex-shrink-0">
                        <svg
                          className="w-4 h-4 sm:w-5 sm:h-5 text-white"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth="2"
                            d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"
                          ></path>
                        </svg>
                      </div>
                      <div>
                        <div className="text-sm font-medium text-gray-900 dark:text-white">
                          Amount to Send
                        </div>
                        <div className="text-xs text-gray-500 dark:text-gray-400">
                          Primary transaction value
                        </div>
                      </div>
                    </div>
                    <div className="text-left sm:text-right">
                      <div className="text-base sm:text-lg font-bold text-gray-900 dark:text-white">
                        {amount} CCC
                      </div>
                    </div>
                  </div>

                  <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center p-3 sm:p-5 bg-yellow-50 dark:bg-yellow-900/20 rounded-xl border border-yellow-200/50 dark:border-yellow-700/50 space-y-3 sm:space-y-0">
                    <div className="flex items-center space-x-3">
                      <div className="w-8 h-8 sm:w-10 sm:h-10 bg-yellow-500 rounded-full flex items-center justify-center flex-shrink-0">
                        <svg
                          className="w-4 h-4 sm:w-5 sm:h-5 text-white"
                          fill="currentColor"
                          viewBox="0 0 20 20"
                        >
                          <path d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z"></path>
                        </svg>
                      </div>
                      <div>
                        <div className="text-sm font-medium text-gray-900 dark:text-white">
                          Network Fee
                        </div>
                        <div className="text-xs text-gray-500 dark:text-gray-400">
                          Processing fee
                        </div>
                      </div>
                    </div>
                    <div className="text-left sm:text-right">
                      <div className="text-base sm:text-lg font-medium text-gray-900 dark:text-white">
                        {fee} CCC
                      </div>
                    </div>
                  </div>

                  <div className="p-4 sm:p-6 bg-gradient-to-r from-primary-50 to-purple-50 dark:from-primary-900/20 dark:to-purple-900/20 rounded-xl border-2 border-primary-200/50 dark:border-primary-700/50">
                    <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center mb-3 sm:mb-4 space-y-3 sm:space-y-0">
                      <div className="flex items-center space-x-3">
                        <div className="w-10 h-10 sm:w-12 sm:h-12 bg-gradient-primary rounded-full flex items-center justify-center flex-shrink-0">
                          <svg
                            className="w-5 h-5 sm:w-6 sm:h-6 text-white"
                            fill="currentColor"
                            viewBox="0 0 20 20"
                          >
                            <path
                              fillRule="evenodd"
                              d="M4 4a2 2 0 00-2 2v4a2 2 0 002 2V6h10a2 2 0 00-2-2H4zm2 6a2 2 0 012-2h8a2 2 0 012 2v4a2 2 0 01-2 2H8a2 2 0 01-2-2v-4zm6 4a2 2 0 100-4 2 2 0 000 4z"
                              clipRule="evenodd"
                            ></path>
                          </svg>
                        </div>
                        <div>
                          <div className="text-base sm:text-lg font-bold text-gray-900 dark:text-white">
                            Total Cost
                          </div>
                          <div className="text-sm text-gray-600 dark:text-gray-400">
                            Amount + Network Fee
                          </div>
                        </div>
                      </div>
                      <div className="text-left sm:text-right">
                        <div className="text-xl sm:text-2xl font-bold text-gray-900 dark:text-white">
                          {amount + fee} CCC
                        </div>
                      </div>
                    </div>
                    <div className="bg-white/50 dark:bg-white/10 p-3 sm:p-4 rounded-lg">
                      <div className="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                        Calculation Breakdown:
                      </div>
                      <div className="text-xs sm:text-sm text-gray-600 dark:text-gray-400 font-mono break-all">
                        Amount: {amount} CCC + Network Fee: {fee} CCC = Total:{' '}
                        {amount + fee} CCC
                      </div>
                    </div>
                  </div>

                  <div className="flex flex-col sm:flex-row sm:justify-between sm:items-center p-3 sm:p-5 bg-green-50 dark:bg-green-900/20 rounded-xl border border-green-200/50 dark:border-green-700/50 space-y-3 sm:space-y-0">
                    <div className="flex items-center space-x-3">
                      <div className="w-8 h-8 sm:w-10 sm:h-10 bg-green-500 rounded-full flex items-center justify-center flex-shrink-0">
                        <svg
                          className="w-4 h-4 sm:w-5 sm:h-5 text-white"
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
                      <div>
                        <div className="text-sm font-medium text-gray-900 dark:text-white">
                          Remaining Balance
                        </div>
                        <div className="text-xs text-gray-500 dark:text-gray-400">
                          After this transaction
                        </div>
                      </div>
                    </div>
                    <div className="text-left sm:text-right">
                      <div className="text-base sm:text-lg font-medium text-green-700 dark:text-green-400">
                        {balance - amount - fee} CCC
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div className="space-y-4 sm:space-y-6">
              <div className="bg-white/70 dark:bg-gray-800/70 rounded-xl sm:rounded-2xl p-4 sm:p-6 border border-gray-200/50 dark:border-gray-600/50">
                <h4 className="text-base sm:text-lg font-semibold text-gray-900 dark:text-white mb-4 sm:mb-6 flex items-center space-x-2">
                  <svg
                    className="w-4 h-4 sm:w-5 sm:h-5 text-primary-600 dark:text-primary-400"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z"></path>
                  </svg>
                  <span>Network Details</span>
                </h4>

                <div className="space-y-3 sm:space-y-4">
                  <div className="p-3 sm:p-4 bg-green-50 dark:bg-green-900/20 rounded-xl border border-green-200/50 dark:border-green-700/50">
                    <div className="text-xs text-gray-500 dark:text-gray-400 mb-2">
                      Network Status
                    </div>
                    <div className="flex items-center space-x-2 mb-1">
                      <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
                      <span className="text-sm font-semibold text-green-700 dark:text-green-400">
                        Operational
                      </span>
                    </div>
                    <div className="text-xs text-gray-600 dark:text-gray-300">
                      All systems running normally
                    </div>
                  </div>

                  <div className="p-3 sm:p-4 bg-blue-50 dark:bg-blue-900/20 rounded-xl border border-blue-200/50 dark:border-blue-700/50">
                    <div className="text-xs text-gray-500 dark:text-gray-400 mb-2">
                      Transaction Type
                    </div>
                    <div className="text-sm font-semibold text-blue-700 dark:text-blue-400 mb-1">
                      Standard Transfer
                    </div>
                    <div className="text-xs text-gray-600 dark:text-gray-300">
                      Simple CCC token transfer
                    </div>
                  </div>

                  <div className="p-3 sm:p-4 bg-white/50 dark:bg-white/10 rounded-xl">
                    <div className="text-xs text-gray-500 dark:text-gray-400 mb-2">
                      Transaction Status
                    </div>
                    <div className="text-sm font-semibold text-gray-900 dark:text-white mb-1">
                      Ready to Send
                    </div>
                    <div className="text-xs text-gray-600 dark:text-gray-300">
                      All details verified
                    </div>
                  </div>
                </div>
              </div>

              <div className="bg-white/70 dark:bg-gray-800/70 rounded-xl sm:rounded-2xl p-4 sm:p-6 border border-gray-200/50 dark:border-gray-600/50">
                <h4 className="text-base sm:text-lg font-semibold text-gray-900 dark:text-white mb-3 sm:mb-4 flex items-center space-x-2">
                  <svg
                    className="w-4 h-4 sm:w-5 sm:h-5 text-green-600 dark:text-green-400"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2 2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
                      clipRule="evenodd"
                    ></path>
                  </svg>
                  <span>Security Status</span>
                </h4>

                <div className="space-y-2 sm:space-y-3">
                  <div className="flex items-center space-x-3 text-sm">
                    <svg
                      className="w-3 h-3 sm:w-4 sm:h-4 text-green-500 flex-shrink-0"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                        clipRule="evenodd"
                      ></path>
                    </svg>
                    <span className="text-xs sm:text-sm text-gray-700 dark:text-gray-300">
                      Address verified &amp; locked
                    </span>
                  </div>
                  <div className="flex items-center space-x-3 text-sm">
                    <svg
                      className="w-3 h-3 sm:w-4 sm:h-4 text-green-500 flex-shrink-0"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                        clipRule="evenodd"
                      ></path>
                    </svg>
                    <span className="text-xs sm:text-sm text-gray-700 dark:text-gray-300">
                      Secure connection active
                    </span>
                  </div>
                  <div className="flex items-center space-x-3 text-sm">
                    <svg
                      className="w-3 h-3 sm:w-4 sm:h-4 text-green-500 flex-shrink-0"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                        clipRule="evenodd"
                      ></path>
                    </svg>
                    <span className="text-xs sm:text-sm text-gray-700 dark:text-gray-300">
                      Transaction encrypted
                    </span>
                  </div>
                  <div className="flex items-center space-x-3 text-sm">
                    <svg
                      className="w-3 h-3 sm:w-4 sm:h-4 text-green-500 flex-shrink-0"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                        clipRule="evenodd"
                      ></path>
                    </svg>
                    <span className="text-xs sm:text-sm text-gray-700 dark:text-gray-300">
                      Sufficient balance confirmed
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {message && (
            <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-xl border border-blue-200/50 dark:border-blue-700/50">
              <div className="text-sm font-medium text-blue-900 dark:text-blue-100 mb-2">
                Transaction Message:
              </div>
              <div className="text-sm text-gray-700 dark:text-gray-300 italic bg-white/50 dark:bg-white/10 p-3 rounded-lg">
                {message}
              </div>
            </div>
          )}

          <div className="p-4 bg-yellow-50 dark:bg-yellow-900/20 rounded-xl border border-yellow-200/50 dark:border-yellow-700/50">
            <div className="flex items-start space-x-3">
              <svg
                className="w-5 h-5 text-yellow-600 dark:text-yellow-400 mt-0.5 flex-shrink-0"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                  clipRule="evenodd"
                ></path>
              </svg>
              <div>
                <div className="text-sm font-medium text-yellow-800 dark:text-yellow-200 mb-1">
                  Important Reminder
                </div>
                <div className="text-sm text-yellow-700 dark:text-yellow-300">
                  Blockchain transactions are <strong>irreversible</strong>.
                  Please verify all details are correct before confirming.
                </div>
              </div>
            </div>
          </div>
        </div>

        <div className="flex flex-col sm:flex-row space-y-3 sm:space-y-0 sm:space-x-4 mt-6 sm:mt-8">
          <button
            type="button"
            onClick={() => actions.closeModal()}
            disabled={isLoading}
            className="disabled:opacity-50 cursor-pointer flex-1 px-4 sm:px-6 py-2.5 sm:py-3 bg-white/70 dark:bg-gray-800/70 border border-gray-200/50 dark:border-gray-600/50 text-gray-900 dark:text-white rounded-xl sm:rounded-2xl font-medium hover:bg-white/90 dark:hover:bg-gray-800/90 transition-all duration-200 flex items-center justify-center space-x-2"
          >
            <svg
              className="w-4 h-4 sm:w-5 sm:h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M11 17l-5-5m0 0l5-5m-5 5h12"
              ></path>
            </svg>
            <span className="text-sm sm:text-base">Back to Edit</span>
          </button>
          <button
            type="button"
            disabled={isLoading}
            onClick={onSendTransaction}
            className="disabled:opacity-50 cursor-pointer flex-1 px-4 sm:px-6 py-2.5 sm:py-3 bg-primary text-white rounded-xl sm:rounded-2xl font-medium hover:opacity-90 transition-all duration-200 shadow-lg hover:shadow-xl flex items-center justify-center space-x-2"
          >
            <svg
              className="w-4 h-4 sm:w-5 sm:h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"
              ></path>
            </svg>
            <span className="text-sm sm:text-base">
              Confirm &amp; Send Transaction
            </span>
            {isLoading && <Loader2 className="w-5 h-5 animate-spin" />}
          </button>
        </div>
      </div>
    </div>
  );
};

export default TransactionDetailModal;
