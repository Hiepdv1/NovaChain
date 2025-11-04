import CopyButton from '@/components/button/copyButton';
import { TransactionItem } from '../../types/transaction';
import {
  FormatFloat,
  FormatTimestamp,
  TruncateHash,
} from '@/shared/utils/format';
import { ArrowRight, Clock, ExternalLink, Layers, Zap } from 'lucide-react';
import StatusBadge from '@/components/status/tx-status';
import { memo } from 'react';
import Link from 'next/link';

interface TransactionCardProps {
  transaction: TransactionItem;
}

const TransactionCard = ({ transaction }: TransactionCardProps) => {
  return (
    <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-200 dark:border-gray-700 p-6 hover:shadow-xl hover:scale-[1.02] hover:border-blue-300 dark:hover:border-blue-600 transition-all duration-300 cursor-pointer group">
      <div className="flex items-start justify-between mb-5">
        <div className="flex items-center gap-3 flex-1 min-w-0">
          <div className="flex-shrink-0 w-12 h-12 bg-gradient-to-br from-blue-500 to-purple-500 rounded-xl flex items-center justify-center shadow-lg">
            <Zap className="w-6 h-6 text-white" />
          </div>
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <span className="text-sm font-mono font-semibold text-gray-900 dark:text-gray-100 truncate">
                {TruncateHash(transaction.TxID, 12, 10)}
              </span>
              <CopyButton text={transaction.TxID} />
            </div>
            <div className="flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400">
              <Clock className="w-3.5 h-3.5" />
              {FormatTimestamp(transaction.CreateAt)}
            </div>
          </div>
        </div>
        <StatusBadge status={'success'} />
      </div>

      <div className="bg-gradient-to-r from-gray-50 to-blue-50 dark:from-gray-900 dark:to-blue-900/20 rounded-xl p-4 mb-4 border border-gray-100 dark:border-gray-700">
        <div className="flex items-center gap-3 mb-3">
          <div className="flex-1">
            <div className="text-xs text-gray-500 dark:text-gray-400 mb-1 font-medium">
              From
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm font-mono text-blue-600 dark:text-blue-400 font-semibold truncate">
                {TruncateHash(transaction.Fromhash.String, 10, 8)}
              </span>
              <CopyButton text={transaction.Fromhash.String} />
            </div>
          </div>
          <div className="flex-shrink-0">
            <div className="w-8 h-8 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full flex items-center justify-center">
              <ArrowRight className="w-4 h-4 text-white" />
            </div>
          </div>
          <div className="flex-1">
            <div className="text-xs text-gray-500 dark:text-gray-400 mb-1 font-medium">
              To
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm font-mono text-purple-600 dark:text-purple-400 font-semibold truncate">
                {TruncateHash(transaction.Tohash.String, 10, 8)}
              </span>
              <CopyButton text={transaction.Tohash.String} />
            </div>
          </div>
        </div>
      </div>

      <div className="grid grid-cols-3 gap-4">
        <div className="bg-gray-50 dark:bg-gray-900 rounded-xl p-3 border border-gray-100 dark:border-gray-700">
          <div className="text-xs text-gray-500 dark:text-gray-400 mb-1 font-medium">
            Value
          </div>
          <div className="text-base font-bold text-gray-900 dark:text-gray-100">
            {FormatFloat(Number(transaction.Amount.String), 8)}
          </div>
        </div>
        <div className="bg-gray-50 dark:bg-gray-900 rounded-xl p-3 border border-gray-100 dark:border-gray-700">
          <div className="text-xs text-gray-500 dark:text-gray-400 mb-1 font-medium">
            Fee
          </div>
          <div className="text-base font-semibold text-gray-700 dark:text-gray-300">
            {FormatFloat(Number(transaction.Fee.String), 8)}
          </div>
        </div>
        <div className="bg-gray-50 dark:bg-gray-900 rounded-xl p-3 border border-gray-100 dark:border-gray-700">
          <div className="text-xs text-gray-500 dark:text-gray-400 mb-1 font-medium flex items-center gap-1">
            <Layers className="w-3 h-3" />
            Block
          </div>
          <Link
            href={`/blocks/${transaction.BID}`}
            className="text-base font-semibold text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 flex items-center gap-1 group-hover:underline"
          >
            {TruncateHash(transaction.BID)}
            <ExternalLink className="w-3 h-3" />
          </Link>
        </div>
      </div>
    </div>
  );
};

export default memo(TransactionCard);
