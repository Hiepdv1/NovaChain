'use client';

import { useCallback } from 'react';
import { useNetworkOverview } from '../hooks/useDashboardQuery';
import NetworkOverviewSkeleton from './NetworkOverviewSkeleton';
import ErrorState from '@/components/errorState';

const OverviewComponent = () => {
  const { data, isLoading, isError, isFetching, refetch, error } =
    useNetworkOverview();

  const onRetry = useCallback(async () => {
    await refetch();
  }, [refetch]);

  if (isLoading || isFetching) {
    return <NetworkOverviewSkeleton />;
  }

  if (!data || isError) {
    return <ErrorState message={error?.message} onRetry={onRetry} />;
  }

  const overview = data.data;

  return (
    <div className="mb-8">
      <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">
        Network Overview
      </h2>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
        <div className="glass-card bg-gradient-glass dark:bg-gradient-glass-dark rounded-2xl p-6 border border-white/20 dark:border-gray-700/50 hover-lift">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">
              Block Height
            </h3>
            <div className="w-8 h-8 bg-blue-100 dark:bg-blue-900/30 rounded-lg flex items-center justify-center">
              <svg
                className="w-4 h-4 text-blue-600 dark:text-blue-400"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z"></path>
              </svg>
            </div>
          </div>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">
            {overview.Chain.BestHeight}
          </p>
          <p className="text-xs text-green-600 dark:text-green-400 mt-1">
            +{overview.Chain.PerHours} blocks/hr
          </p>
        </div>

        <div className="glass-card bg-gradient-glass dark:bg-gradient-glass-dark rounded-2xl p-6 border border-white/20 dark:border-gray-700/50 hover-lift">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">
              Hashrate
            </h3>
            <div className="w-8 h-8 bg-purple-100 dark:bg-purple-900/30 rounded-lg flex items-center justify-center">
              <svg
                className="w-4 h-4 text-purple-600 dark:text-purple-400"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path d="M13 6a3 3 0 11-6 0 3 3 0 016 0zM18 8a2 2 0 11-4 0 2 2 0 014 0zM14 15a4 4 0 00-8 0v3h8v-3z"></path>
              </svg>
            </div>
          </div>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">
            {overview.Hashrate.Value}
          </p>
          <p className="text-xs text-green-600 dark:text-green-400 mt-1">
            +{overview.Hashrate.Per24H} 24h
          </p>
        </div>

        <div className="glass-card bg-gradient-glass dark:bg-gradient-glass-dark rounded-2xl p-6 border border-white/20 dark:border-gray-700/50 hover-lift">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">
              Total Transactions
            </h3>
            <div className="w-8 h-8 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center">
              <svg
                className="w-4 h-4 text-green-600 dark:text-green-400"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path d="M8 5a1 1 0 100 2h5.586l-1.293 1.293a1 1 0 001.414 1.414L15 8.414V14a1 1 0 11-2 0V9.414l-1.293 1.293a1 1 0 01-1.414-1.414L12 7.586H8z"></path>
              </svg>
            </div>
          </div>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">
            {overview.Transaction.Total}
          </p>
          <p className="text-xs text-blue-600 dark:text-blue-400 mt-1">
            +{overview.Transaction.AddedToday} today
          </p>
        </div>

        <div className="glass-card bg-gradient-glass dark:bg-gradient-glass-dark rounded-2xl p-6 border border-white/20 dark:border-gray-700/50 hover-lift">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">
              Pending Txs
            </h3>
            <div className="w-8 h-8 bg-green-100 dark:bg-green-900/30 rounded-lg flex items-center justify-center">
              <svg
                className="w-4 h-4 text-green-600 dark:text-green-400"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path d="M8 5a1 1 0 100 2h5.586l-1.293 1.293a1 1 0 001.414 1.414L15 8.414V14a1 1 0 11-2 0V9.414l-1.293 1.293a1 1 0 01-1.414-1.414L12 7.586H8z"></path>
              </svg>
            </div>
          </div>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">
            {overview.PendingTx.Count}
          </p>
          <p className="text-xs text-gray-600 dark:text-gray-400 mt-1">
            +{overview.PendingTx.AddedToday} today
          </p>
        </div>

        <div className="glass-card bg-gradient-glass dark:bg-gradient-glass-dark rounded-2xl p-6 border border-white/20 dark:border-gray-700/50 hover-lift">
          <div className="flex items-center justify-between mb-2">
            <h3 className="text-sm font-medium text-gray-600 dark:text-gray-400">
              Active Miners
            </h3>
            <div className="w-8 h-8 bg-yellow-100 dark:bg-yellow-900/30 rounded-lg flex items-center justify-center">
              <svg
                className="w-4 h-4 text-yellow-600 dark:text-yellow-400"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z"
                  clipRule="evenodd"
                ></path>
              </svg>
            </div>
          </div>
          <p className="text-2xl font-bold text-gray-900 dark:text-white">
            {overview.ActiveMiners.Count}
          </p>
          <p className="text-xs text-green-600 dark:text-green-400 mt-1">
            +{overview.ActiveMiners.Worker} active today
          </p>
        </div>
      </div>
    </div>
  );
};

export default OverviewComponent;
