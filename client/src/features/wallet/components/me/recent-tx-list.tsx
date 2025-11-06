'use client';

import { Fragment, useCallback, useState } from 'react';
import { useWalletRecentTransactions } from '../../hook/useWalletQuery';
import TransactionItemSkeleton from './loading/transaction-item-loader';
import ErrorState from '@/components/errorState';
import TransactionItem from './transaction-item';

const RecentTransaction = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const { isLoading, isFetching, isError, error, data, refetch } =
    useWalletRecentTransactions(
      {
        limit: 10,
        page: currentPage,
      },
      {
        gcTime: 0,
        staleTime: 0,
        refetchOnReconnect: false,
        refetchOnWindowFocus: false,
      },
    );

  const onRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  if (isLoading || isFetching) {
    return (
      <Fragment>
        {new Array(10).map((_, idx) => {
          return <TransactionItemSkeleton key={idx} />;
        })}
      </Fragment>
    );
  }

  if (isError || !data) {
    return <ErrorState message={error?.message || ''} onRetry={onRetry} />;
  }

  const recentTransactions = data.data;

  console.log(recentTransactions);

  return (
    <div className="p-2">
      {recentTransactions.map((tx) => (
        <TransactionItem key={tx.ID} tx={tx} />
      ))}
    </div>
  );
};

export default RecentTransaction;
