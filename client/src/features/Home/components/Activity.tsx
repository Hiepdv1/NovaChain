'use client';

import {
  FormatFloat,
  FormatHash,
  FormatSize,
  FormatTimestamp,
} from '@/shared/utils/format';
import { useRecentActivity } from '../hooks/useDashboardQuery';
import ActivityLoadingSkeleton from './ActivityLoadingSkeleton';
import { useCallback } from 'react';
import { CACHE_TIME } from '@/shared/constants/ttl';
import ErrorState from '@/components/errorState';

const ActivityComponent = () => {
  const { data, isLoading, isFetching, isError, error, refetch } =
    useRecentActivity({
      retry: false,
      staleTime: 0,
      gcTime: CACHE_TIME,
    });

  const onRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  if (isLoading || isFetching) {
    return <ActivityLoadingSkeleton />;
  }

  if (isError || !data) {
    return <ErrorState message={error?.message} onRetry={onRetry} />;
  }

  const activity = data.data;

  return (
    <div className="mb-8">
      <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">
        Recent Activity
      </h2>
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="glass bg-gradient-glass dark:bg-gradient-glass-dark rounded-2xl p-6 border border-white/20 dark:border-gray-700/50">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white flex items-center">
              <div className="w-2 h-2 bg-green-500 rounded-full mr-2 animate-pulse"></div>
              Latest Blocks
            </h3>
            <button className="text-sm text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 font-medium">
              View All →
            </button>
          </div>
          <div className="space-y-3">
            {activity.Blocks.map((item) => {
              return (
                <div
                  key={item.BID}
                  className="flex items-center justify-between p-3 bg-white/30 dark:bg-white/5 rounded-xl hover:bg-white/40 dark:hover:bg-white/10 transition-colors cursor-pointer"
                >
                  <div className="flex items-center space-x-3">
                    <div className="w-10 h-10 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center">
                      <svg
                        className="w-5 h-5 text-blue-600 dark:text-blue-400"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4z"></path>
                      </svg>
                    </div>
                    <div>
                      <p className="font-semibold text-gray-900 dark:text-white">
                        #{item.Height}
                      </p>
                      <p className="text-sm text-gray-500 dark:text-gray-400">
                        {item.TxCount} txs • {FormatSize(item.Size)}
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-sm font-medium text-gray-900 dark:text-white">
                      {FormatTimestamp(item.Timestamp)}
                    </p>
                    <p className="text-xs text-green-600 dark:text-green-400">
                      {FormatHash(item.BID)}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        <div className="glass bg-gradient-glass dark:bg-gradient-glass-dark rounded-2xl p-6 border border-white/20 dark:border-gray-700/50">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-white flex items-center">
              <div className="w-2 h-2 bg-blue-500 rounded-full mr-2 animate-pulse"></div>
              Latest Transactions
            </h3>
            <button className="text-sm text-primary-600 dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300 font-medium">
              View All →
            </button>
          </div>
          <div className="space-y-3">
            {activity.Txs.map((item) => {
              return (
                <div
                  key={item.ID}
                  className="flex items-center justify-between p-3 bg-white/30 dark:bg-white/5 rounded-xl hover:bg-white/40 dark:hover:bg-white/10 transition-colors cursor-pointer"
                >
                  <div className="flex items-center space-x-3">
                    <div className="w-10 h-10 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center">
                      <svg
                        className="w-5 h-5 text-green-600 dark:text-green-400"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path d="M8 5a1 1 0 100 2h5.586l-1.293 1.293a1 1 0 001.414 1.414L15 8.414V14a1 1 0 11-2 0V9.414l-1.293 1.293a1 1 0 01-1.414-1.414L12 7.586H8z"></path>
                      </svg>
                    </div>
                    <div>
                      <p className="font-semibold text-gray-900 dark:text-white font-mono text-sm">
                        {FormatHash(item.TxID)}
                      </p>
                      <p className="text-sm text-gray-500 dark:text-gray-400">
                        {FormatHash(item.Fromhash.String)} →{' '}
                        {FormatHash(item.Tohash.String)}
                      </p>
                    </div>
                  </div>
                  <div className="text-right">
                    <p className="text-sm font-medium text-gray-900 dark:text-white">
                      {FormatFloat(Number(item.Amount.String))} CCC
                    </p>
                    <p className="text-xs text-green-600 dark:text-green-400">
                      {FormatTimestamp(item.CreateAt)}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
};

export default ActivityComponent;
