import { ChevronLeft, ChevronRight } from 'lucide-react';
import { memo } from 'react';

interface SearchPagination {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

const SearchPagination = ({
  currentPage,
  totalPages,
  onPageChange,
}: SearchPagination) => {
  const getPageNumbers = () => {
    const pages: (string | number)[] = [];
    const showEllipsisStart = currentPage > 3;
    const showEllipsisEnd = currentPage < totalPages - 2;

    if (totalPages <= 7) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      pages.push(1);

      if (showEllipsisStart) {
        pages.push('...');
      }

      for (
        let i = Math.max(2, currentPage - 1);
        i <= Math.min(currentPage + 1, totalPages - 1);
        i++
      ) {
        pages.push(i);
      }

      if (showEllipsisEnd) {
        pages.push('...');
      }

      pages.push(totalPages);
    }

    return pages;
  };

  const onPage = (page: number) => {
    if (page === currentPage) {
      return;
    }

    onPageChange(page);
  };

  return (
    <div className="flex items-center justify-center gap-2 mt-8">
      <button
        onClick={() => onPage(currentPage - 1)}
        disabled={currentPage === 1}
        className="cursor-pointer p-2 rounded-xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-750 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300 hover:scale-105 disabled:hover:scale-100"
      >
        <ChevronLeft className="w-5 h-5 text-gray-700 dark:text-gray-300" />
      </button>

      {getPageNumbers().map((page, index) =>
        typeof page === 'string' ? (
          <span
            key={`ellipsis-${index}`}
            className="px-4 py-2 text-gray-500 dark:text-gray-400"
          >
            ...
          </span>
        ) : (
          <button
            key={page}
            onClick={() => onPage(page)}
            className={`min-w-[40px] cursor-pointer px-4 py-2 rounded-xl font-medium transition-all duration-300 ${
              currentPage === page
                ? 'bg-gradient-to-r from-blue-500 to-purple-600 text-white shadow-lg scale-110'
                : 'bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-750 hover:scale-105'
            }`}
          >
            {page}
          </button>
        ),
      )}

      <button
        onClick={() => onPage(currentPage + 1)}
        disabled={currentPage === totalPages}
        className="cursor-pointer p-2 rounded-xl bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-750 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-300 hover:scale-105 disabled:hover:scale-100"
      >
        <ChevronRight className="w-5 h-5 text-gray-700 dark:text-gray-300" />
      </button>
    </div>
  );
};

export default memo(SearchPagination);
