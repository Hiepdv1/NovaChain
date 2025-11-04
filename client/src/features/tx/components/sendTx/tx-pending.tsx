import { toast } from '@/components/globalToaster';
import type { TransactionPending } from '../../types/transaction';

const statusConfig = {
  icon: (
    <div className="relative w-16 h-16 lg:w-20 lg:h-20 flex items-center justify-center">
      <div className="absolute inset-0 border-3 border-blue-300 dark:border-blue-400 rounded-full animate-spin opacity-30"></div>
      <div
        className="absolute inset-2 border-3 border-indigo-400 dark:border-indigo-300 rounded-full animate-spin opacity-50"
        style={{ animationDirection: 'reverse', animationDuration: '1.5s' }}
      ></div>
      <div className="w-12 h-12 bg-gradient-to-br from-blue-500 to-indigo-600 dark:from-blue-400 dark:to-indigo-500 rounded-full flex items-center justify-center shadow-lg">
        <span className="text-2xl">⏳</span>
      </div>
    </div>
  ),
  label: 'Transaction Pending',
  bgColor: 'bg-blue-50 dark:bg-blue-950/30',
  borderColor: 'border-blue-200 dark:border-blue-800',
  textColor: 'text-blue-900 dark:text-blue-100',
  accentColor: 'text-blue-600 dark:text-blue-400',
  description: 'Your transaction has been submitted successfully',
  subDescription: 'Waiting in the mempool queue for miners',
  notice:
    'Transaction is securely queued in the mempool. In PoW blockchain, transactions are processed sequentially to maintain security.',
};

const TransactionPending = (data: TransactionPending) => {
  const currentStatus = statusConfig;

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    toast.success('Copied to clipboard');
  };

  const truncateAddress = (address: string, start = 6, end = 4) => {
    if (address.length <= start + end) return address;
    return `${address.slice(0, start)}...${address.slice(-end)}`;
  };

  return (
    <div className="max-w-4xl mx-auto p-4">
      <div
        className={`${currentStatus.bgColor} ${currentStatus.borderColor} border-2 rounded-2xl p-6 mb-6 shadow-sm`}
      >
        <div className="flex items-center space-x-6">
          <div className="flex-shrink-0">{currentStatus.icon}</div>

          <div className="flex-1">
            <div className="flex items-center space-x-3 mb-2">
              <h2 className={`text-2xl font-bold ${currentStatus.textColor}`}>
                {currentStatus.label}
              </h2>
              <span
                className={`px-3 py-1 text-xs font-semibold ${currentStatus.accentColor} bg-white dark:bg-gray-800 rounded-full border ${currentStatus.borderColor}`}
              >
                {status.toUpperCase()}
              </span>
            </div>
            <p className={`${currentStatus.textColor} text-sm mb-1 opacity-90`}>
              {currentStatus.description}
            </p>
            <p className={`${currentStatus.textColor} text-sm opacity-70`}>
              {currentStatus.subDescription}
            </p>
          </div>
        </div>
      </div>

      <div className="bg-blue-50 dark:bg-blue-950/20 border border-blue-200 dark:border-blue-800 rounded-xl p-4 mb-6">
        <div className="flex items-start space-x-3">
          <div className="w-8 h-8 bg-blue-100 dark:bg-blue-900/50 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
            <span className="text-blue-600 dark:text-blue-400 text-sm">ℹ️</span>
          </div>
          <p className="text-blue-800 dark:text-blue-200 text-sm leading-relaxed">
            {currentStatus.notice}
          </p>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div className="space-y-4">
          <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-xl p-5 shadow-sm">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wide">
                From Address
              </h3>
              <button
                onClick={() => copyToClipboard(data.Address)}
                className="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
                title="Copy address"
              >
                <svg
                  className="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                  />
                </svg>
              </button>
            </div>
            <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-3 font-mono text-sm text-gray-700 dark:text-gray-300 break-all">
              <span className="hidden sm:block">{data.Address}</span>
              <span className="sm:hidden">{truncateAddress(data.Address)}</span>
            </div>
          </div>

          <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-xl p-5 shadow-sm">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wide">
                To Address
              </h3>
              <button
                onClick={() => copyToClipboard(data.ReceiverAddress)}
                className="p-2 text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors rounded-lg hover:bg-gray-100 dark:hover:bg-gray-800"
                title="Copy address"
              >
                <svg
                  className="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"
                  />
                </svg>
              </button>
            </div>
            <div className="bg-gray-50 dark:bg-gray-800 rounded-lg p-3 font-mono text-sm text-gray-700 dark:text-gray-300 break-all">
              <span className="hidden sm:block">{data.ReceiverAddress}</span>
              <span className="sm:hidden">
                {truncateAddress(data.ReceiverAddress)}
              </span>
            </div>
          </div>
        </div>

        <div className="space-y-4">
          <div className="bg-green-50 dark:bg-green-950/20 border border-green-200 dark:border-green-800 rounded-xl p-5 shadow-sm">
            <h3 className="text-sm font-semibold text-green-700 dark:text-green-300 uppercase tracking-wide mb-3">
              Amount
            </h3>
            <div className="flex items-center space-x-2">
              <span className="text-2xl font-bold text-green-800 dark:text-green-200">
                {data.Amount}
              </span>
              <span className="text-lg font-medium text-green-600 dark:text-green-400">
                CCC
              </span>
            </div>
          </div>

          <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-xl p-5 shadow-sm">
            <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wide mb-3">
              Network Fee
            </h3>
            <div className="flex items-center space-x-2">
              <span className="text-xl font-bold text-gray-800 dark:text-gray-200">
                {data.Fee}
              </span>
              <span className="text-base font-medium text-gray-600 dark:text-gray-400">
                CCC
              </span>
            </div>
          </div>

          <div className="bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-700 rounded-xl p-5 shadow-sm">
            <div className="flex items-center justify-between mb-3">
              <h3 className="text-sm font-semibold text-gray-700 dark:text-gray-300 uppercase tracking-wide">
                Transaction Hash
              </h3>
              <button className="flex items-center space-x-1 text-xs text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 transition-colors px-3 py-1 bg-blue-50 dark:bg-blue-950/30 rounded-full">
                <svg
                  className="w-3 h-3"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14"
                  />
                </svg>
                <span>View</span>
              </button>
            </div>
            <div
              onClick={() => copyToClipboard(data.TxID)}
              className="bg-gray-50 dark:bg-gray-800 rounded-lg p-3 font-mono text-sm text-blue-600 dark:text-blue-400 break-all cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors"
              title="Click to copy"
            >
              <span className="hidden lg:block">{data.TxID}</span>
              <span className="lg:hidden">
                {truncateAddress(data.TxID, 8, 6)}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default TransactionPending;
