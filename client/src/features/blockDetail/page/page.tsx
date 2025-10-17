'use client';

import { useParams, useSearchParams } from 'next/navigation';
import { useBlockDetailQuery } from '../hook/useBlockQuery';
import { FormatSize, FormatTimestamp, IsNumber } from '@/shared/utils/format';
import NotFoundState from '@/components/notfound';
import BlockDetailSkeletonLoader from '../components/Block-detail-loading';
import {
  Activity,
  Check,
  ChevronLeft,
  Clock,
  Copy,
  Database,
  HardDrive,
  Hash,
  User,
  Zap,
} from 'lucide-react';
import { useCallback, useState } from 'react';
import ErrorState from '@/components/errorState';
import TransactionList from '../components/Transaction-list';
import { useHistoryContext } from '@/components/providers/History-provider';

const BlockDetailPage = () => {
  const param = useParams();
  const query = useSearchParams();
  const { b_hash } = param;
  const queryPage = query.get('page');
  let currentPage = 1;

  if (IsNumber(queryPage) && Number(queryPage) > 0) {
    currentPage = Number(Number(queryPage).toFixed());
  }
  const [copiedHash, setCopiedHash] = useState<string | null>(null);
  const { history, back } = useHistoryContext();

  const { isError, isLoading, isFetched, data, error, refetch } =
    useBlockDetailQuery(
      {
        b_hash: b_hash as string,
        limit: 2,
        page: currentPage,
      },
      {
        enabled: !!b_hash && b_hash.length === 64,
        retry: false,
        refetchOnWindowFocus: false,
        refetchOnReconnect: false,
        staleTime: 0,
        gcTime: 0,
      },
    );

  const copyToClipboard = (text: string, id: string) => {
    navigator.clipboard.writeText(text);
    setCopiedHash(id);
    setTimeout(() => setCopiedHash(null), 2000);
  };

  const onRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  if (isLoading || !isFetched) {
    return <BlockDetailSkeletonLoader />;
  }

  if (isError) {
    if (error.statusCode === 404) {
      return <NotFoundState />;
    }

    return <ErrorState onRetry={onRetry} />;
  }

  if (!data) {
    return <NotFoundState />;
  }

  const blockData = data.data;

  return (
    <div>
      {history.length > 0 && (
        <div className="sticky top-0 z-10 backdrop-blur-md border-b bg-white/90 dark:bg-slate-800/90 border-slate-200 dark:border-slate-700 transition-colors rounded-2xl">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
            <button
              onClick={() => back()}
              className="flex items-center gap-2 px-4 py-2 rounded-lg transition-all hover:scale-105 hover:bg-blue-100 dark:hover:bg-slate-700 text-slate-900 dark:text-slate-100"
            >
              <ChevronLeft className="w-5 h-5" />
              <span className="font-medium">Back to Blocks</span>
            </button>
          </div>
        </div>
      )}

      <div className="max-w-7xl mx-auto py-8">
        <div className="rounded-xl p-6 mb-6 transition-all hover:shadow-xl bg-gradient-to-br from-white to-blue-50 dark:from-slate-800 dark:to-slate-700 shadow-lg shadow-blue-200/50 dark:shadow-blue-500/10">
          <div className="flex items-start justify-between mb-4">
            <div>
              <h1 className="text-2xl font-bold mb-2 text-slate-900 dark:text-slate-50">
                Block #{blockData.Height.toLocaleString()}
              </h1>
              <div className="flex items-center gap-2 text-sm text-slate-600 dark:text-slate-400">
                <Clock className="w-4 h-4" />
                <span>{FormatTimestamp(blockData.Timestamp)}</span>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-3 p-4 rounded-lg bg-gradient-to-r from-blue-50 to-purple-50 dark:from-slate-700 dark:to-slate-800 border border-slate-200 dark:border-slate-700">
            <Hash className="w-5 h-5 text-blue-500 dark:text-blue-400" />
            <code
              className={`flex-1 text-sm font-mono break-all text-slate-700 dark:text-slate-300`}
            >
              {blockData.BID}
            </code>
            <button
              onClick={() => copyToClipboard(blockData.BID, 'block')}
              className="cursor-pointer p-2 rounded-lg transition-all hover:scale-110 hover:bg-blue-100 dark:hover:bg-slate-700"
            >
              {copiedHash === 'block' ? (
                <Check className="w-5 h-5 text-green-500 dark:text-green-400" />
              ) : (
                <Copy className="w-5 h-5 text-slate-600 dark:text-slate-400" />
              )}
            </button>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-8">
        {[
          {
            label: 'Block Height',
            value: `#${blockData.Height.toLocaleString()}`,
            icon: Activity,
          },
          { label: 'Transactions', value: blockData.TxCount, icon: Zap },
          { label: 'Total Fees', value: blockData.TotalFee, icon: Database },
          {
            label: 'Size',
            value: FormatSize(blockData.Size),
            icon: HardDrive,
          },
          { label: 'Difficulty', value: blockData.Difficulty, icon: Activity },
          { label: 'Miner', value: blockData.Miner, icon: User },
        ].map((item, idx) => {
          const Icon = item.icon;

          return (
            <div
              key={idx}
              className="rounded-xl p-5 transition-all hover:shadow-lg hover:scale-[1.02] bg-gradient-to-br from-white to-blue-50 dark:from-slate-800 dark:to-slate-700 shadow-md shadow-blue-200/30 dark:shadow-blue-500/10"
            >
              <div className="flex items-center gap-3 mb-2">
                <div className="p-2 rounded-lg bg-gradient-to-br from-blue-400 to-purple-400">
                  <Icon className="w-5 h-5 text-white dark:text-black/40" />
                </div>
                <span className="text-sm font-medium text-slate-600 dark:text-slate-400">
                  {item.label}
                </span>
              </div>
              <p className="text-sm truncate font-semibold font-mono text-slate-900 dark:text-slate-100">
                {item.value}
              </p>
            </div>
          );
        })}
      </div>

      <TransactionList blockData={data.data} />
    </div>
  );
};

export default BlockDetailPage;
