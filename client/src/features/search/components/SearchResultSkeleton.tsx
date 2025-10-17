'use client';

export default function SearchResultSkeleton() {
  return (
    <div className="animate-pulse">
      {/* Header */}
      <div className="mb-6 flex items-center justify-between">
        <div>
          <div className="h-6 w-40 bg-gray-300 dark:bg-gray-700 rounded mb-2" />
          <div className="h-4 w-64 bg-gray-200 dark:bg-gray-800 rounded" />
        </div>
        <div className="h-4 w-24 bg-gray-200 dark:bg-gray-800 rounded" />
      </div>

      {/* List items */}
      <div className="space-y-3 sm:space-y-4">
        {Array.from({ length: 5 }).map((_, i) => (
          <div
            key={i}
            className="p-4 sm:p-5 bg-white/60 dark:bg-white/10 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm"
          >
            <div className="flex items-center justify-between">
              <div className="space-y-2">
                <div className="h-4 w-48 bg-gray-300 dark:bg-gray-700 rounded" />
                <div className="h-3 w-32 bg-gray-200 dark:bg-gray-800 rounded" />
              </div>
              <div className="h-4 w-20 bg-gray-200 dark:bg-gray-800 rounded" />
            </div>
          </div>
        ))}
      </div>

      {/* Pagination */}
      <div className="flex items-center justify-center gap-2 mt-8">
        <div className="h-8 w-8 bg-gray-200 dark:bg-gray-800 rounded-lg" />
        <div className="h-8 w-8 bg-gray-300 dark:bg-gray-700 rounded-lg" />
        <div className="h-8 w-8 bg-gray-200 dark:bg-gray-800 rounded-lg" />
      </div>
    </div>
  );
}
