'use client';

const TransactionListSkeleton = () => {
  return (
    <div className="space-y-4">
      {[1, 2, 3, 4, 5].map((i) => (
        <div
          key={i}
          className="rounded-lg p-4 border animate-pulse bg-gradient-to-br from-gray-100 to-gray-200 dark:from-slate-900/50 dark:to-slate-800 border-gray-200 dark:border-slate-700"
        >
          <div className="flex items-start justify-between mb-3">
            <div className="flex items-center gap-2 flex-1 min-w-0">
              <div className="w-4 h-4 rounded bg-gray-300 dark:bg-slate-700 flex-shrink-0" />
              <div className="h-4 bg-gray-300 dark:bg-slate-700 rounded w-32" />
              <div className="w-5 h-5 rounded bg-gray-300 dark:bg-slate-700" />
            </div>
            <div className="h-3 bg-gray-300 dark:bg-slate-700 rounded w-20" />
          </div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-sm">
            <div>
              <div className="h-3 bg-gray-300 dark:bg-slate-700 rounded w-16 mb-2" />
              <div className="h-3 bg-gray-300 dark:bg-slate-700 rounded w-40" />
            </div>
            <div>
              <div className="h-3 bg-gray-300 dark:bg-slate-700 rounded w-12 mb-2" />
              <div className="h-3 bg-gray-300 dark:bg-slate-700 rounded w-36" />
            </div>
          </div>

          <div className="flex items-center justify-between mt-3 pt-3 border-t border-gray-200 dark:border-slate-700">
            <div className="flex items-center gap-4 w-full sm:w-auto">
              <div>
                <div className="h-3 bg-gray-300 dark:bg-slate-700 rounded w-10 mb-1" />
                <div className="h-4 bg-gray-300 dark:bg-slate-700 rounded w-16" />
              </div>
              <div>
                <div className="h-3 bg-gray-300 dark:bg-slate-700 rounded w-8 mb-1" />
                <div className="h-4 bg-gray-300 dark:bg-slate-700 rounded w-14" />
              </div>
            </div>

            <div className="hidden sm:flex items-center gap-1 px-6 py-2 rounded-lg text-xs bg-gray-300 dark:bg-slate-700 w-28 h-7" />
          </div>
        </div>
      ))}
    </div>
  );
};

export default TransactionListSkeleton;
