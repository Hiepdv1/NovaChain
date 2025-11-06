const BalanceCardSkeleton = () => (
  <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 p-5 sm:p-6 lg:p-8 animate-pulse">
    <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-6">
      <div className="flex-1">
        <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-24 mb-3"></div>
        <div className="h-12 bg-gray-200 dark:bg-gray-700 rounded w-48 mb-2"></div>
        <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-32"></div>
      </div>
      <div className="flex gap-3">
        <div className="h-12 bg-gray-200 dark:bg-gray-700 rounded-lg w-28"></div>
        <div className="h-12 bg-gray-200 dark:bg-gray-700 rounded-lg w-28"></div>
      </div>
    </div>
    <div className="mt-6 pt-6 border-t border-gray-200 dark:border-gray-700">
      <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-full"></div>
    </div>
  </div>
);

export default BalanceCardSkeleton;
