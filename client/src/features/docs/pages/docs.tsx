'use client';

import NovaChain from '../components/nova-chain-page';
import { FormatFloat } from '@/shared/utils/format';
import { useDownloadInfoQuery } from '../hooks/useDownloadQuery';
import ErrorState from '@/components/errorState';
import { useCallback } from 'react';

const DocumentPage = () => {
  const { isLoading, isFetching, isError, data, error, refetch } =
    useDownloadInfoQuery();

  const onRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  if (isLoading || isFetching) {
    return (
      <div className="h-screen flex justify-center py-10">
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

  const info = data.data;

  const sizeFile = FormatFloat(info.Size / 1024 / 1024, 2);

  return <NovaChain SizeFile={sizeFile} />;
};

export default DocumentPage;
