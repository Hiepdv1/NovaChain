const StatCardSkeleton = () => (
  <div className="bg-white dark:bg-gray-800 rounded-xl p-5 border border-gray-200 dark:border-gray-700 animate-pulse">
    <div className="w-12 h-12 bg-gray-200 dark:bg-gray-700 rounded-lg mb-3"></div>
    <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-20 mb-2"></div>
    <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded w-16 mb-1"></div>
    <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-12"></div>
  </div>
);

export default StatCardSkeleton;
