import { useState } from 'react';
import { TransactionPending } from '../../types/transaction';
import { ChevronDown, ChevronUp, Clock, Loader2, XCircle } from 'lucide-react';
import { FormatFloat, TruncateHash } from '@/shared/utils/format';
import CopyButton from '@/components/button/copyButton';
import TxPendingStatusBadge from './pending-tx-status';

interface PendingTransactionCardProps {
  transaction: TransactionPending;
}

const TimeBadge = ({ time }: { time: string }) => {
  return (
    <div className="flex items-center gap-2 px-3 py-1.5 bg-blue-50 dark:bg-blue-900/20 rounded-lg border border-blue-200 dark:border-blue-800">
      <Clock className="w-4 h-4 text-blue-600 dark:text-blue-400" />
      <span className="text-sm font-semibold text-blue-700 dark:text-blue-300">
        {new Date(time).toLocaleDateString()}
      </span>
    </div>
  );
};

const PendingTransactionCard = ({
  transaction,
}: PendingTransactionCardProps) => {
  const [expanded, setExpanded] = useState(false);

  const getStatusColor = (status: 'pending' | 'failed') => {
    if (status === 'pending') {
      return 'border-amber-300 dark:border-amber-700 bg-amber-50/50 dark:bg-amber-900/10';
    }
    if (status === 'failed') {
      return 'border-red-300 dark:border-red-700 bg-red-50/50 dark:bg-red-900/10';
    }
    return 'border-gray-200 dark:border-gray-700';
  };

  const getStatusIcon = (status: 'pending' | 'failed') => {
    if (status === 'failed') return XCircle;
    return Loader2;
  };

  const getStatusGradient = (status: 'pending' | 'failed') => {
    if (status === 'failed') return 'from-red-500 to-pink-500';
    return 'from-amber-500 to-orange-500';
  };

  const StatusIcon = getStatusIcon(transaction.Status);

  return (
    <div
      className={`bg-white dark:bg-gray-800 rounded-2xl shadow-sm border-2 ${getStatusColor(
        transaction.Status,
      )} p-6 hover:shadow-xl transition-all duration-300`}
    >
      {/* Header */}
      <div className="flex items-start justify-between mb-5">
        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-3 mb-2">
            <div
              className={`flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center shadow-lg bg-gradient-to-br ${getStatusGradient(
                transaction.Status,
              )}`}
            >
              <StatusIcon
                className={`w-6 h-6 text-white ${
                  transaction.Status !== 'failed' ? 'animate-spin' : ''
                }`}
              />
            </div>
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 mb-1">
                <span className="text-sm font-mono font-bold text-gray-900 dark:text-gray-100 truncate">
                  {TruncateHash(transaction.TxID, 12, 10)}
                </span>
                <CopyButton type="icon" text={transaction.TxID} />
              </div>
              <TimeBadge time={transaction.CreatedAt.Time} />
            </div>
          </div>
        </div>
        <TxPendingStatusBadge
          status={transaction.Status}
          fee={FormatFloat(Number(transaction.Fee), 8)}
        />
      </div>

      {/* Main Info */}
      <div className="bg-gradient-to-r from-gray-50 to-amber-50 dark:from-gray-900 dark:to-amber-900/20 rounded-xl p-4 mb-4 border border-gray-100 dark:border-gray-700">
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <div className="text-xs text-gray-500 dark:text-gray-400 mb-1 font-medium">
              From
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm font-mono text-amber-600 dark:text-amber-400 font-semibold truncate">
                {TruncateHash(transaction.Address, 10, 8)}
              </span>
              <CopyButton text={transaction.Address} />
            </div>
          </div>
          <div>
            <div className="text-xs text-gray-500 dark:text-gray-400 mb-1 font-medium">
              To
            </div>
            <div className="flex items-center gap-2">
              <span className="text-sm font-mono text-orange-600 dark:text-orange-400 font-semibold truncate">
                {TruncateHash(transaction.ReceiverAddress, 10, 8)}
              </span>
              <CopyButton text={transaction.ReceiverAddress} />
            </div>
          </div>
        </div>
      </div>

      {/* Quick Stats */}
      <div className="grid grid-cols-2 gap-3 mb-4">
        <div className="bg-gray-50 dark:bg-gray-900 rounded-xl p-3 border border-gray-100 dark:border-gray-700">
          <div className="text-xs text-gray-500 dark:text-gray-400 mb-1 font-medium">
            Amount
          </div>
          <div className="text-base font-bold text-gray-900 dark:text-gray-100">
            {FormatFloat(Number(transaction.Amount), 8)}
          </div>
        </div>
        <div className="bg-gray-50 dark:bg-gray-900 rounded-xl p-3 border border-gray-100 dark:border-gray-700">
          <div className="text-xs text-gray-500 dark:text-gray-400 mb-1 font-medium">
            Transaction Fee
          </div>
          <div className="text-base font-semibold text-gray-700 dark:text-gray-300">
            {FormatFloat(Number(transaction.Fee), 8)}
          </div>
        </div>
      </div>

      {/* Status Info */}
      {transaction.Status === 'pending' && (
        <div className="mb-4 p-3 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg">
          <div className="flex items-center gap-2 text-sm text-amber-700 dark:text-amber-300">
            <Loader2 className="w-4 h-4 animate-spin" />
            <span className="font-medium">
              Waiting in mempool to be mined into a block
            </span>
          </div>
        </div>
      )}

      {transaction.Status === 'failed' && transaction.Message.Valid && (
        <div className="mb-4 p-3 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg">
          <div className="text-xs font-medium text-red-600 dark:text-red-400 mb-1">
            Failure Reason
          </div>
          <div className="text-sm text-red-700 dark:text-red-300">
            {transaction.Message.String}
          </div>
        </div>
      )}

      <button
        onClick={() => setExpanded(!expanded)}
        className="w-full flex items-center justify-center gap-2 py-2 text-sm font-semibold text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200 transition-colors"
      >
        {expanded ? (
          <>
            <span>Show Less</span>
            <ChevronUp className="w-4 h-4" />
          </>
        ) : (
          <>
            <span>Show Details</span>
            <ChevronDown className="w-4 h-4" />
          </>
        )}
      </button>

      {expanded && (
        <div className="mt-4 pt-4 border-t border-gray-200 dark:border-gray-700 space-y-3">
          <div className="grid grid-cols-2 gap-4">
            <div>
              <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">
                Transaction ID
              </div>
              <div className="text-sm font-mono font-semibold text-gray-900 dark:text-gray-100 break-all">
                {transaction.TxID}
              </div>
            </div>
            <div>
              <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">
                Priority Level
              </div>
              <div className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                {transaction.Priority.Valid
                  ? transaction.Priority.Int32 === 1
                    ? 'High'
                    : transaction.Priority.Int32 === 2
                    ? 'Medium'
                    : 'Low'
                  : 'N/A'}
              </div>
            </div>
          </div>

          <div>
            <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">
              Submitted At
            </div>
            <div className="text-sm font-semibold text-gray-900 dark:text-gray-100">
              {new Date(transaction.CreatedAt.Time).toLocaleDateString()}
            </div>
          </div>

          {transaction.UpdatedAt.Valid && (
            <div>
              <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">
                Last Updated
              </div>
              <div className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                {new Date(transaction.UpdatedAt.Time).toLocaleDateString()}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default PendingTransactionCard;
