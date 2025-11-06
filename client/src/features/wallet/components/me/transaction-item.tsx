import { ArrowDownLeft, ArrowUpRight } from 'lucide-react';
import { RecentTransaction } from '../../types/wallet';
import { FormatFloat, FormatTimestamp } from '@/shared/utils/format';
import { memo } from 'react';

interface TransactionItemProps {
  tx: RecentTransaction;
}

const TransactionItem = ({ tx }: TransactionItemProps) => {
  const isReceive = tx.Type === 'received';

  return (
    <div className="flex items-center gap-3 sm:gap-4 p-3 sm:p-4 hover:bg-gray-50 dark:hover:bg-gray-700/50 rounded-lg transition-all cursor-pointer">
      <div
        className={`p-2.5 rounded-lg flex-shrink-0 ${
          isReceive
            ? 'bg-green-100 dark:bg-green-900/30'
            : 'bg-blue-100 dark:bg-blue-900/30'
        }`}
      >
        {isReceive ? (
          <ArrowDownLeft className="w-5 h-5 text-green-600 dark:text-green-400" />
        ) : (
          <ArrowUpRight className="w-5 h-5 text-blue-600 dark:text-blue-400" />
        )}
      </div>

      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 mb-1">
          <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">
            {isReceive ? 'Received from' : 'Sent to'}
          </span>
        </div>
        <div className="text-xs text-gray-500 dark:text-gray-400 font-mono truncate">
          {tx.Fromhash.String}
        </div>
        <div className="text-xs text-gray-400 dark:text-gray-500 mt-1">
          {FormatTimestamp(tx.CreateAt)}
        </div>
      </div>

      <div className="text-right flex-shrink-0">
        <div
          className={`text-sm sm:text-base font-bold ${
            isReceive
              ? 'text-green-600 dark:text-green-400'
              : 'text-rose-600 dark:text-rose-400'
          }`}
        >
          {isReceive ? '+' : '-'}
          {FormatFloat(Number(tx.Amount.String), 8)}
        </div>
        <div
          className={`text-xs px-2 py-0.5 rounded mt-1 inline-block bg-green-100 dark:bg-green-900/50 text-green-700 dark:text-green-400`}
        >
          confirmed
        </div>
      </div>
    </div>
  );
};

export default memo(TransactionItem);
