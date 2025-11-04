const PendingTransactionsLoader = () => {
  return (
    <div className="animate-pulse space-y-4">
      {[...Array(4)].map((_, idx) => (
        <div
          key={idx}
          className="bg-white dark:bg-gray-800 rounded-2xl border-2 border-gray-200 dark:border-gray-700 p-6"
        >
          <div className="flex items-center justify-between mb-4">
            <div className="h-6 bg-gray-200 dark:bg-gray-700 rounded w-1/4"></div>
            <div className="h-8 bg-gray-200 dark:bg-gray-700 rounded-full w-24"></div>
          </div>
          <div className="space-y-3">
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-full"></div>
            <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-3/4"></div>
          </div>
        </div>
      ))}
    </div>
  );
};

export default PendingTransactionsLoader;
