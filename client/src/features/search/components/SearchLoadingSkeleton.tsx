const SearchLoadingSkeleton = ({}) => {
  return (
    <div className="space-y-4">
      {[1, 2, 3, 4, 5].map((item) => (
        <div
          key={item}
          className="p-6 bg-gradient-to-br from-gray-200 to-gray-300 dark:from-gray-800 dark:to-gray-700 rounded-2xl animate-pulse"
        >
          <div className="flex items-center gap-3 mb-4">
            <div className="w-12 h-12 bg-gray-300 dark:bg-gray-600 rounded-xl"></div>
            <div className="flex-1">
              <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-24 mb-2"></div>
              <div className="h-3 bg-gray-300 dark:bg-gray-600 rounded w-32"></div>
            </div>
          </div>
          <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-full mb-2"></div>
          <div className="h-4 bg-gray-300 dark:bg-gray-600 rounded w-2/3"></div>
        </div>
      ))}
    </div>
  );
};

export default SearchLoadingSkeleton;
