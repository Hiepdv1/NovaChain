'use client';

import { useRouter, useSearchParams } from 'next/navigation';
import { useSearchQuery } from '../hooks/useSeachQuery';
import SearchEmptyState from '../components/SearchEmpty';
import UnifiedResultItem from '../components/UnifiedResultItem';
import SearchPagination from '../components/SearchPagination';
import { useCallback } from 'react';
import { CACHE_TIME } from '@/shared/constants/ttl';
import { IsNumber } from '@/shared/utils/format';
import ErrorState from '@/components/errorState';
import SearchResultSkeleton from '../components/SearchResultSkeleton';
import EmptyState from '../components/SearchEmpty';

const SearchResultPage = () => {
  const searchParams = useSearchParams();
  const router = useRouter();
  const query = searchParams.get('result_search');
  const limit = 10;
  const pageQuery = searchParams.get('page');
  let currentPage = 1;

  if (IsNumber(pageQuery)) {
    currentPage = Number(Number(pageQuery).toFixed());
  }

  const { isError, isLoading, isFetching, data, error, refetch } =
    useSearchQuery(
      {
        limit,
        page: currentPage,
        search_query: query || '',
      },
      {
        enabled: !!query && query.trim() !== '' && query.length > 1,
        retry: false,
        refetchOnWindowFocus: false,
        refetchOnReconnect: false,
        staleTime: 0,
        gcTime: CACHE_TIME,
      },
    );

  const onPageChange = useCallback(
    (page: number) => {
      router.push(`/search?result_search=${query}&page=${page}`, {
        scroll: true,
      });
    },
    [router, query],
  );

  const onRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  if (!query || query.length < 1) {
    return <EmptyState />;
  }

  if (isLoading || isFetching) {
    return <SearchResultSkeleton />;
  }

  if (isError || !data) {
    return <ErrorState message={error?.message} onRetry={onRetry} />;
  }

  const results = data.data;

  if (!results || results.length === 0) {
    return <SearchEmptyState />;
  }

  const meta = data.meta;

  return (
    <div>
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-2">
            Search Results
          </h2>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            Found {results.length} results • Showing{' '}
            {(meta.currentPage - 1) * limit || (results.length > 0 ? 1 : 0)}-
            {(meta.currentPage - 1) * meta.limit + results.length} of{' '}
            {meta.total}
          </p>
        </div>
        <div className="text-sm text-gray-500 dark:text-gray-400">
          Page {meta.currentPage} of {Math.ceil(meta.total / meta.limit)}
        </div>
      </div>

      <div className="space-y-3 sm:space-y-4">
        {results.map((item, index) => (
          <UnifiedResultItem key={index} item={item} index={index} />
        ))}
      </div>

      <SearchPagination
        currentPage={meta.currentPage}
        totalPages={Math.ceil(meta.total / meta.limit)}
        onPageChange={onPageChange}
      />
    </div>
  );
};

export default SearchResultPage;
