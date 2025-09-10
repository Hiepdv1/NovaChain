import { formatAddress } from '@/lib/utils';

interface TransactionPreviewProps {
  fromAddr: string;
  recipient?: string;
  active: boolean;
  data: {
    amount: number;
    fee: number;
    message: string;
  };
}

const TransactionPreview = ({
  fromAddr,
  recipient,
  active,
  data,
}: TransactionPreviewProps) => {
  return (
    <div className="sticky top-24 glass-card dark:bg-primary-dark dark:border-secondary-dark rounded-2xl p-6 shadow-lg shadow-black/5 dark:shadow-black/20">
      <div className="flex items-center justify-between mb-4">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
          Transaction Preview
        </h3>
        <div className="flex items-center space-x-2">
          <div className="w-2 h-2 bg-blue-500 rounded-full animate-pulse"></div>
          <span className="font-semibold text-xs text-gray-500 dark:text-gray-400">
            Live Preview
          </span>
        </div>
      </div>

      <div className={`${active ? 'opacity-100' : 'opacity-60'}`}>
        <div className="flex items-center justify-between p-4 bg-gradient-to-r from-blue-50 to-purple-50 dark:from-blue-900/20 dark:to-purple-900/20 rounded-xl">
          <div className="flex flex-col items-center">
            <div className="w-12 h-12 bg-primary rounded-full flex items-center justify-center mb-2">
              <svg
                className="w-6 h-6 text-white"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path d="M4 4a2 2 0 00-2 2v1h16V6a2 2 0 00-2-2H4zM18 9H2v5a2 2 0 002 2h12a2 2 0 002-2V9zM4 13a1 1 0 011-1h1a1 1 0 110 2H5a1 1 0 01-1-1zm5-1a1 1 0 100 2h1a1 1 0 100-2H9z"></path>
              </svg>
            </div>

            <div className="text-xs font-medium text-gray-900 dark:text-white">
              Your wallet
            </div>

            <div className="text-xs text-gray-500 dark:text-gray-400 font-mono">
              {formatAddress(fromAddr)}
            </div>
          </div>

          <div className="flex-1 flex items-center justify-center">
            <div className="flex items-center space-x-2">
              <div className="w-8 h-0.5 bg-primary rounded-full"></div>
              <svg
                className="w-5 h-5 text-primary-600 dark:text-primary-400"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"
                ></path>
              </svg>
              <div className="w-8 h-0.5 bg-gradient-primary rounded-full"></div>
            </div>
          </div>

          <div className="flex flex-col items-center">
            <div className="w-12 h-12 bg-gray-200 dark:bg-gray-700 rounded-full flex items-center justify-center mb-2">
              <svg
                className="w-6 h-6 text-gray-500 dark:text-gray-400"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z"></path>
              </svg>
            </div>

            <div className="text-xs font-medium text-gray-900 dark:text-white">
              Recipient
            </div>

            <div className="text-xs text-gray-500 dark:text-gray-400 font-mono">
              {recipient ? formatAddress(recipient) : 'Not Set'}
            </div>
          </div>
        </div>

        <div className="mt-3 space-y-3">
          <div className="flex justify-between items-center p-3 bg-white/30 dark:bg-white/5 rounded-xl">
            <span className="text-sm text-gray-600 dark:text-gray-400">
              Amount:
            </span>
            <div className="text-right">
              <div className="text-sm font-semibold text-gray-900 dark:text-white">
                {data.amount > 0 ? data.amount : '-'} CCC
              </div>
            </div>
          </div>

          <div className="flex justify-between items-center p-3 bg-white/30 dark:bg-white/5 rounded-xl">
            <span className="text-sm text-gray-600 dark:text-gray-400">
              Network Fee:
            </span>
            <div className="text-right">
              <div className="text-sm text-gray-900 dark:text-white">
                {data.fee > 0 ? data.fee : '-'} CCC
              </div>
            </div>
          </div>

          <div className="p-3 bg-gradient-to-r from-primary-50 to-purple-50 dark:from-primary-900/20 dark:to-purple-900/20 rounded-xl border border-primary-200/50 dark:border-primary-700/50">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm font-semibold text-gray-900 dark:text-white">
                Total Cost (Amount + Fee):
              </span>
              <div className="text-right">
                <div
                  id="previewTotal"
                  className="text-sm font-bold text-gray-900 dark:text-white"
                >
                  {data.amount > 0 ? data.amount : '-'} CCC
                </div>
              </div>
            </div>
            <div className="text-xs text-gray-600 dark:text-gray-400 bg-white/30 dark:bg-white/10 p-2 rounded-lg">
              <span id="totalBreakdown">
                Amount: {data.amount > 0 ? data.amount : '-'} CCC + Mining Fee:{' '}
                {data.fee > 0 ? data.fee : '-'} CCC = Total:{' '}
                {data.amount > 0 && data.fee > 0 ? data.amount + data.fee : '-'}{' '}
                CCC
              </span>
            </div>
          </div>

          <div className="text-center p-3 bg-white/30 dark:bg-white/5 rounded-xl">
            <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">
              Transaction Status
            </div>
            <div className="text-sm font-medium text-gray-900 dark:text-white">
              Ready to Send
            </div>
          </div>

          {data.message && (
            <div className="p-3 bg-blue-50 dark:bg-blue-900/20 rounded-xl border border-blue-200/50 dark:border-blue-700/50">
              <div className="text-xs text-blue-600 dark:text-blue-400 mb-1">
                Transaction Message:
              </div>
              <div className="text-sm text-gray-900 dark:text-white italic break-all whitespace-normal">
                {data.message}
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default TransactionPreview;
