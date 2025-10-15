'use client';

const ActivityLoadingSkeleton = () => {
  return (
    <div className="mb-8 animate-pulse">
      <h2 className="text-2xl font-bold text-gray-900 dark:text-white mb-6">
        Recent Activity
      </h2>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {[1, 2].map((col) => (
          <div
            key={col}
            className="glass bg-gradient-glass dark:bg-gradient-glass-dark rounded-2xl p-6 border border-white/20 dark:border-gray-700/50"
          >
            <div className="h-5 bg-gray-300 dark:bg-gray-700 rounded w-1/3 mb-4"></div>
            <div className="space-y-3">
              {[...Array(5)].map((_, i) => (
                <div
                  key={i}
                  className="flex items-center justify-between p-3 bg-white/30 dark:bg-white/5 rounded-xl"
                >
                  <div className="flex items-center space-x-3 w-full">
                    <div className="w-10 h-10 bg-gray-300 dark:bg-gray-600 rounded-lg"></div>
                    <div className="flex-1">
                      <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-24 mb-2"></div>
                      <div className="h-3 bg-gray-300 dark:bg-gray-600 rounded w-32"></div>
                    </div>
                    <div className="w-16 h-3 bg-gray-300 dark:bg-gray-600 rounded"></div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default ActivityLoadingSkeleton;
