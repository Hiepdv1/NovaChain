import { ArrowRightLeft, Box, ChevronRight, Clock, Copy } from 'lucide-react';
import { SearchItem } from '../types/search';
import { memo, useState } from 'react';
import { FormatSize, FormatTimestamp } from '@/shared/utils/format';
import Link from 'next/link';
import { useHistoryContext } from '@/components/providers/History-provider';
import useCurrentUrl from '@/shared/hooks/useContextUrl';

const UnifiedResultItem = ({
  index,
  item,
}: {
  item: SearchItem;
  index: number;
}) => {
  const [copied, setCopied] = useState(false);

  const { push } = useHistoryContext();
  const currentPath = useCurrentUrl();

  const handleCopy = (text: string) => {
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  const onNavigate = () => {
    push({
      path: currentPath,
      title: 'Back to search results',
    });
  };

  const getTypeConfig = () => {
    switch (item.Type) {
      case 'block':
        return {
          icon: Box,
          label: 'Block',
          gradient: 'from-blue-500 to-blue-600',
          borderColor: 'hover:border-blue-400 dark:hover:border-blue-500',
          bgGlow: 'from-blue-500/10 to-purple-500/10',
          shadowColor:
            'group-hover:shadow-blue-500/30 dark:group-hover:shadow-blue-500/50',
          labelColor: 'text-blue-600 dark:text-blue-400',
          labelBg:
            'bg-blue-100 dark:bg-blue-500/20 text-blue-700 dark:text-blue-300',
        };
      case 'transaction':
        return {
          icon: ArrowRightLeft,
          label: 'Transaction',
          gradient: 'from-purple-500 to-purple-600',
          borderColor: 'hover:border-purple-400 dark:hover:border-purple-500',
          bgGlow: 'from-purple-500/10 to-pink-500/10',
          shadowColor:
            'group-hover:shadow-purple-500/30 dark:group-hover:shadow-purple-500/50',
          labelColor: 'text-purple-600 dark:text-purple-400',
          labelBg:
            'bg-purple-100 dark:bg-purple-500/20 text-purple-700 dark:text-purple-300',
        };
      default:
        return {
          icon: Box,
          label: 'Unknown',
          gradient: 'from-gray-400 to-gray-500',
          borderColor: 'hover:border-gray-400 dark:hover:border-gray-500',
          bgGlow: 'from-gray-400/10 to-gray-500/10',
          shadowColor:
            'group-hover:shadow-gray-500/30 dark:group-hover:shadow-gray-500/50',
          labelColor: 'text-gray-600 dark:text-gray-400',
          labelBg:
            'bg-gray-100 dark:bg-gray-500/20 text-gray-700 dark:text-gray-300',
        };
    }
  };

  const config = getTypeConfig();
  const Icon = config.icon;

  return (
    <Link
      onClick={onNavigate}
      href={
        item.Type === 'block'
          ? `/blocks/${item.Keyword}`
          : `/txs/${item.Keyword}`
      }
      className="block"
    >
      <div
        style={{ animationDelay: `${index * 50}ms` }}
        className={`group animate-fadeInUp relative bg-gradient-to-br from-white via-white to-gray-50/50 dark:from-gray-800 dark:via-gray-800 dark:to-gray-900/50 rounded-2xl p-4 sm:p-6 border border-gray-200/60 dark:border-gray-700/60 hover:shadow-2xl hover:scale-[1.02] ${config.borderColor} transition-all duration-500 cursor-pointer overflow-hidden backdrop-blur-sm`}
      >
        {/* Glow effect */}
        <div
          className={`absolute top-0 right-0 w-24 h-24 sm:w-32 sm:h-32 bg-gradient-to-br ${config.bgGlow} rounded-full blur-3xl group-hover:scale-150 transition-transform duration-700 opacity-50 dark:opacity-100`}
        ></div>

        <div className="relative">
          {/* Header section */}
          <div className="flex flex-col sm:flex-row items-start sm:items-start justify-between mb-4 gap-3">
            <div className="flex items-start gap-3 flex-1 min-w-0">
              <div
                className={`flex-shrink-0 p-2.5 sm:p-3 bg-gradient-to-br ${config.gradient} rounded-xl shadow-lg ${config.shadowColor} transition-shadow duration-300`}
              >
                <Icon className="w-5 h-5 sm:w-6 sm:h-6 text-white" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 flex-wrap">
                  <span
                    className={`text-xs font-bold ${config.labelColor} uppercase tracking-wider`}
                  >
                    {config.label}
                  </span>
                  {item.Type === 'block' && (
                    <span
                      className={`px-2 py-0.5 text-xs font-semibold ${config.labelBg} rounded-full`}
                    >
                      #{item.Data.height.toLocaleString()}
                    </span>
                  )}
                  {item.Type === 'transaction' && (
                    <span className="px-2 py-0.5 text-xs font-semibold bg-green-100 dark:bg-green-500/20 text-green-700 dark:text-green-400 rounded-full border border-green-200 dark:border-green-500/30">
                      Confirmed
                    </span>
                  )}
                </div>

                {item.Data.timestamp && (
                  <div className="flex items-center gap-2 mt-2 text-xs text-gray-600 dark:text-gray-400">
                    <Clock className="w-3 h-3 flex-shrink-0" />
                    <span className="truncate">
                      {FormatTimestamp(item.Data.timestamp)}
                    </span>
                  </div>
                )}
              </div>
            </div>

            <button
              onClick={(e) => {
                e.stopPropagation();
                handleCopy(item.Keyword);
              }}
              className="flex-shrink-0 p-2 cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700/50 rounded-lg transition-colors active:scale-95"
            >
              {copied ? (
                <span className="text-xs text-green-600 dark:text-green-400 font-medium whitespace-nowrap">
                  Copied!
                </span>
              ) : (
                <Copy className="w-4 h-4 text-gray-400 dark:text-gray-500 hover:text-gray-600 dark:hover:text-gray-300" />
              )}
            </button>
          </div>

          {/* Hash/ID section */}
          <div className="mb-4">
            <p className="text-xs font-mono text-gray-500 dark:text-gray-500 mb-1.5">
              {item.Type === 'block'
                ? 'Hash'
                : item.Type === 'transaction'
                ? 'Transaction ID'
                : ''}
            </p>
            <p className="font-mono text-xs sm:text-sm text-gray-900 dark:text-gray-200 break-all leading-relaxed bg-gray-50 dark:bg-gray-900/50 p-3 rounded-lg border border-gray-100 dark:border-gray-700/50">
              {item.Keyword}
            </p>
          </div>

          {/* Block details */}
          {item.Type === 'block' && (
            <div className="grid grid-cols-3 gap-3 sm:gap-4 pt-4 border-t border-gray-200 dark:border-gray-700/70">
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-500 mb-1 truncate">
                  Transactions
                </p>
                <p className="text-sm font-semibold text-gray-900 dark:text-gray-100 truncate">
                  {item.Data.tx_count}
                </p>
              </div>
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-500 mb-1 truncate">
                  Size
                </p>
                <p className="text-sm font-semibold text-gray-900 dark:text-gray-100 truncate">
                  {FormatSize(item.Data.size)}
                </p>
              </div>
              <div>
                <p className="text-xs text-gray-500 dark:text-gray-500 mb-1 truncate">
                  Miner
                </p>
                <p
                  className="text-sm font-semibold text-gray-900 dark:text-gray-100 truncate"
                  title={item.Data.miner}
                >
                  {item.Data.miner}
                </p>
              </div>
            </div>
          )}

          {/* Transaction details */}
          {item.Type === 'transaction' && (
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 pt-4 border-t border-gray-200 dark:border-gray-700/70">
              <div className="bg-gradient-to-br from-gray-50 to-gray-100/50 dark:from-gray-900/50 dark:to-gray-800/50 p-3 rounded-lg border border-gray-100 dark:border-gray-700/50">
                <p className="text-xs text-gray-500 dark:text-gray-500 mb-1.5">
                  Amount
                </p>
                <p className="text-base sm:text-lg font-bold text-gray-900 dark:text-gray-100 truncate">
                  {item.Data.amount} CCC
                </p>
              </div>
              <div className="bg-gradient-to-br from-gray-50 to-gray-100/50 dark:from-gray-900/50 dark:to-gray-800/50 p-3 rounded-lg border border-gray-100 dark:border-gray-700/50">
                <p className="text-xs text-gray-500 dark:text-gray-500 mb-1.5">
                  Transaction Fee
                </p>
                <p className="text-sm font-semibold text-gray-900 dark:text-gray-100 truncate">
                  {item.Data.fee} CCC
                </p>
              </div>
            </div>
          )}

          {/* Chevron indicator */}
          <div className="absolute bottom-3 sm:bottom-4 right-3 sm:right-4 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
            <ChevronRight
              className={`w-4 h-4 sm:w-5 sm:h-5 ${config.labelColor}`}
            />
          </div>
        </div>
      </div>
    </Link>
  );
};

export default memo(UnifiedResultItem);
