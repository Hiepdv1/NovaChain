const TransactionItemSkeleton = () => (
  <div className="flex items-center gap-3 sm:gap-4 p-3 sm:p-4 rounded-lg animate-pulse">
    <div className="w-12 h-12 bg-gray-200 dark:bg-gray-700 rounded-lg flex-shrink-0"></div>
    <div className="flex-1 min-w-0 space-y-2">
      <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-24"></div>
      <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-full max-w-xs"></div>
      <div className="h-3 bg-gray-200 dark:bg-gray-700 rounded w-20"></div>
    </div>
    <div className="text-right flex-shrink-0 space-y-2">
      <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-20 ml-auto"></div>
      <div className="h-5 bg-gray-200 dark:bg-gray-700 rounded w-16 ml-auto"></div>
    </div>
  </div>
);

export default TransactionItemSkeleton;
