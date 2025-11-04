import { ChevronLeft, ChevronRight } from 'lucide-react';
import { memo } from 'react';

interface BlockListPagination {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

const Pagination = ({
  currentPage,
  totalPages,
  onPageChange,
}: BlockListPagination) => {
  return (
    <div className="flex items-center justify-center gap-2">
      <button
        onClick={() => onPageChange(Math.max(1, currentPage - 1))}
        disabled={currentPage === 1}
        className={`cursor-pointer p-2 rounded-lg transition-all duration-300 hover:bg-gray-200 dark:hover:bg-gray-800 ${
          currentPage === 1 ? 'opacity-50 !cursor-not-allowed' : ''
        }`}
      >
        <ChevronLeft size={20} />
      </button>

      {[...Array(totalPages)].map((_, i) => {
        const page = i + 1;
        if (
          page === 1 ||
          page === totalPages ||
          (page >= currentPage - 1 && page <= currentPage + 1)
        ) {
          return (
            <button
              key={page}
              onClick={() => onPageChange(page)}
              className={`cursor-pointer px-4 py-2 rounded-lg transition-all duration-300 ${
                currentPage === page
                  ? 'bg-blue-500 text-white shadow-lg scale-110'
                  : 'hover:bg-gray-200 dark:hover:bg-gray-800'
              }`}
            >
              {page}
            </button>
          );
        } else if (page === currentPage - 2 || page === currentPage + 2) {
          return <span key={page}>...</span>;
        }
        return null;
      })}

      <button
        onClick={() => onPageChange(Math.min(totalPages, currentPage + 1))}
        disabled={currentPage === totalPages}
        className={`cursor-pointer p-2 rounded-lg transition-all duration-300 hover:bg-gray-200 dark:hover:bg-gray-800 ${
          currentPage === totalPages ? 'opacity-50 !cursor-not-allowed' : ''
        }`}
      >
        <ChevronRight size={20} />
      </button>
    </div>
  );
};

export default memo(Pagination);
