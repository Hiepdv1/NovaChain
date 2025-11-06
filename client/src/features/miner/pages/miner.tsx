'use client';

import { useCallback, useState } from 'react';
import { useMiners } from '../hook/useMinerQuery';
import MinerSkeletonLoader from '../components/miner-loader';
import ErrorState from '@/components/errorState';
import { Cpu } from 'lucide-react';
import MinerCard from '../components/miner-card';
import Pagination from '@/components/pagination';

const MinerPage = () => {
  const [currentPage, setCurrentPage] = useState(1);
  const { isLoading, isFetching, isError, error, data, refetch } = useMiners(
    {
      limit: 10,
      page: currentPage,
    },
    {
      refetchOnWindowFocus: false,
      refetchOnReconnect: false,
      retry: false,
      gcTime: 0,
      staleTime: 0,
    },
  );

  const onRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  if (isLoading || isFetching) {
    return <MinerSkeletonLoader />;
  }

  if (isError || !data) {
    return <ErrorState message={error?.message || ''} onRetry={onRetry} />;
  }

  const miners = data.data;
  const meta = data.meta;

  if (miners.length === 0) {
    return (
      <div className="min-h-screen">
        <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-200 dark:border-gray-700 p-12 text-center">
          <Cpu className="w-16 h-16 text-gray-400 dark:text-gray-600 mx-auto mb-4" />
          <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-2">
            No Miners Found
          </h3>
          <p className="text-gray-600 dark:text-gray-400">
            No miners match your search criteria.
          </p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen">
      <div className="space-y-6">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {miners.map((miner, idx) => (
            <MinerCard key={miner.MinerPubkey} miner={miner} rank={idx + 1} />
          ))}
        </div>

        {meta.total > meta.limit && (
          <Pagination
            currentPage={meta.currentPage}
            onPageChange={setCurrentPage}
            totalPages={Math.ceil(meta.total / meta.limit)}
          />
        )}
      </div>
    </div>
  );
};

export default MinerPage;
