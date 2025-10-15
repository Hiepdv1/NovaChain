import { FileSearch } from 'lucide-react';

const EmptyState = () => {
  return (
    <div className="flex flex-col items-center justify-center py-20 px-4">
      <div className="relative">
        <div className="absolute inset-0 bg-gradient-to-r from-blue-500 to-purple-500 rounded-full blur-3xl opacity-20 animate-pulse"></div>
        <div className="relative p-8 bg-gradient-to-br from-gray-100 to-gray-200 dark:from-gray-800 dark:to-gray-700 rounded-full mb-6">
          <FileSearch className="w-16 h-16 text-gray-400 dark:text-gray-500" />
        </div>
      </div>
      <h3 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-3">
        No results found
      </h3>
      <p className="text-gray-600 dark:text-gray-400 text-center max-w-md leading-relaxed">
        We couldn&apos;t find any blocks, transactions, or outputs matching your
        search. Try a different query.
      </p>
    </div>
  );
};

export default EmptyState;
