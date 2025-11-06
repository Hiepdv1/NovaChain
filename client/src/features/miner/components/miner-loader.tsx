const MinerSkeletonLoader = () => (
  <div className="animate-pulse space-y-4">
    {[...Array(6)].map((_, idx) => (
      <div
        key={idx}
        className="bg-white dark:bg-gray-800 rounded-2xl border border-gray-200 dark:border-gray-700 p-6"
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-4 flex-1">
            <div className="w-16 h-16 bg-gray-200 dark:bg-gray-700 rounded-xl"></div>
            <div className="flex-1 space-y-2">
              <div className="h-5 bg-gray-200 dark:bg-gray-700 rounded w-1/3"></div>
              <div className="h-4 bg-gray-200 dark:bg-gray-700 rounded w-1/2"></div>
            </div>
          </div>
          <div className="h-10 bg-gray-200 dark:bg-gray-700 rounded w-24"></div>
        </div>
      </div>
    ))}
  </div>
);

export default MinerSkeletonLoader;
