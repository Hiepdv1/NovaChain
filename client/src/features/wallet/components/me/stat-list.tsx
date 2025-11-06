'use client';

import { Fragment, useCallback } from 'react';
import { useWalletSummary } from '../../hook/useWalletQuery';
import StatCardSkeleton from './loading/stat-card-loader';
import ErrorState from '@/components/errorState';
import { ArrowDownLeft, ArrowUpRight, Layers } from 'lucide-react';
import StatCard, { StatCardProps } from './stat-card';
import { FormatFloat } from '@/shared/utils/format';

const StatList = () => {
  const { isLoading, isFetching, isError, error, data, refetch } =
    useWalletSummary({
      gcTime: 0,
      staleTime: 0,
      refetchOnReconnect: false,
      refetchOnWindowFocus: false,
      retry: false,
    });

  const onRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  if (isLoading || isFetching) {
    return (
      <Fragment>
        <StatCardSkeleton />
        <StatCardSkeleton />
        <StatCardSkeleton />
      </Fragment>
    );
  }

  if (isError || !data) {
    return <ErrorState message={error?.message || ''} onRetry={onRetry} />;
  }

  const statData = data.data;

  const stats: StatCardProps[] = [
    {
      Icon: Layers,
      label: 'Total Transactions',
      value: statData.TotalTx,
      color: 'purple',
    },
    {
      Icon: ArrowUpRight,
      label: 'Total Sent',
      value: FormatFloat(Number(statData.TotalSent), 8),
      subValue: 'CCC',
      color: 'blue',
    },
    {
      Icon: ArrowDownLeft,
      label: 'Total Received',
      value: FormatFloat(Number(statData.TotalReceived), 8),
      subValue: 'CCC',
      color: 'green',
    },
  ];

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-2 sm:gap-3">
      {stats.map((stat, idx) => (
        <StatCard
          key={idx}
          Icon={stat.Icon}
          color={stat.color}
          label={stat.label}
          subValue={stat.subValue}
          value={stat.value}
        />
      ))}
    </div>
  );
};

export default StatList;
