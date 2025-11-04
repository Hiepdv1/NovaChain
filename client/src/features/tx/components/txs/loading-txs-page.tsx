const TransactionSkeletonLoader = () => (
  <div className="animate-pulse space-y-4">
    {[...Array(5)].map((_, idx) => (
      <div
        key={idx}
        className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-200 dark:border-gray-700 p-6"
      >
        <div className="flex items-center justify-between mb-4">
          <div className="h-5 bg-gray-200 dark:bg-gray-700 rounded w-1/3"></div>
          <div className="h-6 bg-gray-200 dark:bg-gray-700 rounded-full w-20"></div>
        </div>
        <div className="space-y-3">
          <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-full"></div>
          <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-5/6"></div>
          <div className="flex gap-4">
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/4"></div>
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/4"></div>
          </div>
        </div>
      </div>
    ))}
  </div>
);

export default TransactionSkeletonLoader;
