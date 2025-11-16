'use client';

import { useParams } from 'next/navigation';
import { useTransactionDetail } from '../hook/useTransactionQuery';
import ErrorState from '@/components/errorState';
import { useCallback } from 'react';
import TransactionCard from '../components/tx-detail/transaction-card';

const TransactionDetailPage = ({}) => {
  const { tx_hash } = useParams();

  const { isLoading, isFetching, isError, data, error, refetch } =
    useTransactionDetail(tx_hash as string);

  const onRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  if (isLoading || isFetching) {
    return (
      <div className="flex justify-center py-10">
        <div
          className="
        h-6 w-6 
        border-4 border-gray-300 dark:border-gray-700 
        border-t-blue-600 dark:border-t-blue-400
        rounded-full animate-spin
      "
        />
      </div>
    );
  }

  if (isError || !data) {
    return <ErrorState message={error?.message} onRetry={onRetry} />;
  }

  const transaction = data.data;

  return <TransactionCard transaction={transaction} />;
};

export default TransactionDetailPage;
