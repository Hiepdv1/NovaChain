'use client';

import ErrorState from '@/components/errorState';
import ListBlockLoadingSkeleton from '../components/blocks/Loading-list-block';
import { useListBlocks } from '../hooks/useBlockQuery';
import BlockCard from '../components/blocks/BlockCard';
import BlockListPagination from '../components/blocks/Block-list-pagination';
import { useCallback, useState } from 'react';

const BlockPage = () => {
  const [currentPage, setCurrentPage] = useState(1);

  const { isLoading, isFetching, data, error, isError, refetch } =
    useListBlocks(
      {
        limit: 12,
        page: currentPage,
      },
      {
        retry: false,
        refetchOnWindowFocus: false,
        refetchOnReconnect: false,
        staleTime: 0,
        gcTime: 0,
      },
    );

  const onRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  if (isFetching || isLoading) {
    return <ListBlockLoadingSkeleton />;
  }

  if (isError || !data) {
    return <ErrorState message={error?.message || ''} onRetry={onRetry} />;
  }

  return (
    <div className="min-h-screen transition-colors duration-300 text-gray-900 dark:text-gray-100">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        {data.data.map((block) => (
          <BlockCard key={block.BID} block={block} />
        ))}
      </div>

      {data.meta.total > data.meta.limit && (
        <BlockListPagination
          currentPage={currentPage}
          totalPages={Math.ceil(data.meta.total / data.meta.limit)}
          onPageChange={setCurrentPage}
        />
      )}
    </div>
  );
};

export default BlockPage;
