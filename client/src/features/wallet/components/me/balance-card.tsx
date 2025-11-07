'use client';

import { memo, useCallback, useState } from 'react';
import BalanceCardSkeleton from './loading/balance-loader';
import ErrorState from '@/components/errorState';
import { useWalletQuery } from '../../hook/useWalletQuery';
import { Eye, EyeOff, Send } from 'lucide-react';
import { FormatFloat } from '@/shared/utils/format';
import CopyButton from '@/components/button/copyButton';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

const BalanceCard = () => {
  const [showBalance, setShowBalance] = useState(true);
  const { isLoading, isFetching, isError, error, data, refetch } =
    useWalletQuery();
  const router = useRouter();

  const onRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  if (isFetching || isLoading) {
    return <BalanceCardSkeleton />;
  }

  if (isError || !data) {
    if (error && error.statusCode === 401) {
      router.push('/');
      return null;
    }

    return <ErrorState message={error?.message || ''} onRetry={onRetry} />;
  }

  const walletData = data.data;

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-5 sm:p-6 lg:p-8 transition-colors">
      <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-2">
            <span className="text-sm font-medium text-gray-600 dark:text-gray-400">
              Total Balance
            </span>
            <button
              onClick={() => setShowBalance(!showBalance)}
              className="cursor-pointer p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700 transition-all"
            >
              {showBalance ? (
                <Eye className="w-4 h-4 text-gray-400" />
              ) : (
                <EyeOff className="w-4 h-4 text-gray-400" />
              )}
            </button>
          </div>

          {showBalance ? (
            <>
              <div className="flex items-baseline flex-wrap gap-2 mb-2">
                <span className="text-3xl sm:text-2xl lg:text-xl font-bold text-gray-900 dark:text-gray-100">
                  {FormatFloat(Number(walletData.Balance), 8)}
                </span>
                <span className="text-lg sm:text-xl text-gray-600 dark:text-gray-400">
                  CCC
                </span>
              </div>
            </>
          ) : (
            <div className="text-3xl sm:text-2xl lg:text-xl font-bold text-gray-900 dark:text-gray-100">
              ••••••••
            </div>
          )}
        </div>

        <Link href="/txs/send" className="flex cursor-pointer gap-3">
          <button className="flex-1 lg:flex-none flex items-center justify-center gap-2 px-5 sm:px-6 py-3 bg-gradient-to-r from-blue-600 to-blue-700 hover:from-blue-700 hover:to-blue-800 text-white rounded-lg font-semibold transition-all shadow-md hover:shadow-lg">
            <Send className="w-5 h-5" />
            <span>Send</span>
          </button>
        </Link>
      </div>

      <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
        <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
          <span className="text-sm text-gray-600 dark:text-gray-400">
            Wallet Address
          </span>
          <div className="flex items-center gap-2">
            <span className="text-sm font-mono text-gray-900 dark:text-gray-100 truncate max-w-[200px] sm:max-w-none">
              {walletData.Address.String}
            </span>
            <CopyButton text={walletData.Address.String} />
          </div>
        </div>
      </div>
    </div>
  );
};

export default memo(BalanceCard);
