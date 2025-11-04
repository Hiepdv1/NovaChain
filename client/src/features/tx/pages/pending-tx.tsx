'use client';

import { useCallback, useState } from 'react';
import { usePendingTransactions } from '../hook/useTransactionQuery';
import ErrorState from '@/components/errorState';
import PendingTransactionsLoader from '../components/pending-tx/pending-tx-loader';
import PendingTransactionCard from '../components/pending-tx/tx-card';
import { Inbox } from 'lucide-react';
import Pagination from '@/components/pagination';

const PendingTransactionsPage = () => {
  const [currentPage, setCurrentPage] = useState(1);

  const { isLoading, isFetching, isError, error, data, refetch } =
    usePendingTransactions({
      page: currentPage,
      limit: 10,
    });

  const onRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  if (isLoading || isFetching) {
    return <PendingTransactionsLoader />;
  }

  if (isError || !data) {
    return <ErrorState message={error?.message || ''} onRetry={onRetry} />;
  }

  const meta = data.meta;
  const txPendings = data.data || [];

  if (txPendings.length === 0) {
    return (
      <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border-2 border-gray-200 dark:border-gray-700 p-12 text-center">
        <Inbox className="w-16 h-16 text-gray-400 dark:text-gray-600 mx-auto mb-4" />
        <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-2">
          No Transactions Found
        </h3>
        <p className="text-gray-600 dark:text-gray-400">
          There are currently no transactions in the mempool matching your
          filters.
        </p>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
        {txPendings.map((item) => {
          return <PendingTransactionCard key={item.ID} transaction={item} />;
        })}
      </div>

      {meta.total > meta.limit && (
        <Pagination
          currentPage={meta.currentPage}
          onPageChange={setCurrentPage}
          totalPages={Math.ceil(meta.total / meta.limit)}
        />
      )}
    </div>
  );
};

export default PendingTransactionsPage;
