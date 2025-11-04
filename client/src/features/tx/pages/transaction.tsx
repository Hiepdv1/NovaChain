'use client';

import { useCallback, useState } from 'react';
import { useListTransactions } from '../hook/useTransactionQuery';
import TransactionSkeletonLoader from '../components/txs/loading-txs-page';
import ErrorState from '@/components/errorState';
import TransactionCard from '../components/txs/tx-card';
import Pagination from '@/components/pagination';

const TransactionPage = () => {
  const [currentPage, setCurrentPage] = useState(1);

  const { isLoading, isFetching, isError, error, data, refetch } =
    useListTransactions({
      limit: 10,
      page: currentPage,
    });

  const onRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  if (isFetching || isLoading) {
    return <TransactionSkeletonLoader />;
  }

  if (isError || !data) {
    return <ErrorState onRetry={onRetry} message={error?.message || ''} />;
  }

  return (
    <div className="min-h-screen">
      <div className="space-y-6 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
        {data.data.map((item) => {
          return <TransactionCard key={item.ID} transaction={item} />;
        })}

        {data.meta.total > data.meta.limit && (
          <Pagination
            currentPage={data.meta.currentPage}
            totalPages={Math.ceil(data.meta.total / data.meta.limit)}
            onPageChange={setCurrentPage}
          />
        )}
      </div>
    </div>
  );
};

export default TransactionPage;
