const NetworkStatus = () => {
  return (
    <div className="glass-card dark:bg-primary-dark dark:border-secondary-dark rounded-2xl p-6 shadow-lg shadow-black/5 dark:shadow-black/20">
      <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-4 flex items-center space-x-2">
        <svg
          className="w-5 h-5 text-primary-600 dark:text-primary-400"
          fill="currentColor"
          viewBox="0 0 20 20"
        >
          <path d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z"></path>
        </svg>
        <span>Network Status</span>
      </h3>

      <div className="space-y-4">
        <div className="flex items-center justify-between p-3 bg-green-50 dark:bg-green-900/20 rounded-xl">
          <div className="flex items-center space-x-3">
            <div className="w-3 h-3 bg-green-500 rounded-full animate-pulse"></div>
            <span className="text-sm font-medium text-gray-900 dark:text-white">
              Network Health
            </span>
          </div>
          <span className="text-sm font-semibold text-green-600 dark:text-green-400">
            Excellent
          </span>
        </div>

        <div className="grid grid-cols-2 gap-3">
          <div className="text-center p-3 bg-white/30 dark:bg-white/5 rounded-xl">
            <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">
              Block Height
            </div>
            <div className="text-sm font-semibold text-gray-900 dark:text-white">
              #2,847,392
            </div>
          </div>
          <div className="text-center p-3 bg-white/30 dark:bg-white/5 rounded-xl">
            <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">
              Avg Block Time
            </div>
            <div className="text-sm font-semibold text-gray-900 dark:text-white">
              12.3s
            </div>
          </div>
          <div className="text-center p-3 bg-white/30 dark:bg-white/5 rounded-xl">
            <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">
              Pending Txs
            </div>
            <div className="text-sm font-semibold text-gray-900 dark:text-white">
              1,847
            </div>
          </div>
          <div className="text-center p-3 bg-white/30 dark:bg-white/5 rounded-xl">
            <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">
              Network Load
            </div>
            <div className="text-sm font-semibold text-green-600 dark:text-green-400">
              Normal
            </div>
          </div>
        </div>

        <div className="p-4 bg-blue-50 dark:bg-blue-900/20 rounded-xl border border-blue-200/50 dark:border-blue-700/50">
          <div className="flex items-center space-x-2 mb-2">
            <svg
              className="w-4 h-4 text-blue-600 dark:text-blue-400"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
                clipRule="evenodd"
              ></path>
            </svg>
            <span className="text-sm font-medium text-blue-900 dark:text-blue-100">
              Mining Fee Recommendation
            </span>
          </div>
          <p className="text-xs text-blue-800 dark:text-blue-200">
            Current mining difficulty is normal. Standard mining fee (0.0025
            CCC) provides good balance between cost and mining priority. Higher
            fees increase your transaction&apos;s priority in the mining pool.
          </p>
        </div>
      </div>
    </div>
  );
};

export default NetworkStatus;
