import { ChevronRight, History, Send, Settings } from 'lucide-react';
import BalanceCard from '../components/me/balance-card';
import StatList from '../components/me/stat-list';
import RecentTransaction from '../components/me/recent-tx-list';

const WalletAccessPage = () => {
  return (
    <div className="min-h-screen">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
        <div className="space-y-4 sm:space-y-6">
          <BalanceCard />
          <StatList />

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-4 sm:gap-6">
            <div className="lg:col-span-2">
              <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 transition-colors">
                <div className="px-4 sm:px-6 py-4 sm:py-5 border-b border-gray-200 dark:border-gray-700">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-2 sm:gap-3">
                      <div className="p-2 rounded-lg bg-gradient-to-br from-purple-500 to-purple-600">
                        <History className="w-4 h-4 sm:w-5 sm:h-5 text-white" />
                      </div>
                      <h2 className="text-base sm:text-lg font-bold text-gray-900 dark:text-gray-100">
                        Recent Transactions
                      </h2>
                    </div>
                    <button className="flex items-center gap-1 text-xs sm:text-sm font-medium text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 transition-colors">
                      <span className="hidden sm:inline">View All</span>
                      <span className="sm:hidden">All</span>
                      <ChevronRight className="w-4 h-4" />
                    </button>
                  </div>
                </div>

                <RecentTransaction />
              </div>
            </div>

            <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-5 sm:p-6 transition-colors">
              <h3 className="text-base sm:text-lg font-bold text-gray-900 dark:text-gray-100 mb-4">
                Quick Actions
              </h3>
              <div className="space-y-2">
                <button className="w-full flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-all text-left group">
                  <div className="flex items-center gap-3">
                    <div className="p-1.5 rounded-lg bg-blue-100 dark:bg-blue-900/30">
                      <Send className="w-4 h-4 text-blue-600 dark:text-blue-400" />
                    </div>
                    <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
                      Send Payment
                    </span>
                  </div>
                  <ChevronRight className="w-4 h-4 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300 transition-colors" />
                </button>
                <button className="w-full flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-all text-left group">
                  <div className="flex items-center gap-3">
                    <div className="p-1.5 rounded-lg bg-purple-100 dark:bg-purple-900/30">
                      <History className="w-4 h-4 text-purple-600 dark:text-purple-400" />
                    </div>
                    <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
                      Transaction History
                    </span>
                  </div>
                  <ChevronRight className="w-4 h-4 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300 transition-colors" />
                </button>
                <button className="w-full flex items-center justify-between p-3 rounded-lg hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-all text-left group">
                  <div className="flex items-center gap-3">
                    <div className="p-1.5 rounded-lg bg-gray-100 dark:bg-gray-700">
                      <Settings className="w-4 h-4 text-gray-600 dark:text-gray-400" />
                    </div>
                    <span className="text-sm font-medium text-gray-900 dark:text-gray-100">
                      Wallet Settings
                    </span>
                  </div>
                  <ChevronRight className="w-4 h-4 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300 transition-colors" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default WalletAccessPage;
