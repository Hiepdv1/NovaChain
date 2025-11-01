import { ChevronLeft, ChevronRight } from 'lucide-react';
import { Fragment } from 'react';

interface TransactionPaginationProps {
  totalPages: number;
  currentPage: number;
  goToPage: (page: number) => void;
}

const TransactionPagination = ({
  currentPage,
  totalPages,
  goToPage,
}: TransactionPaginationProps) => {
  const getPageNumbers = () => {
    const pages = [];
    const maxPagesToShow = 5;

    if (totalPages <= maxPagesToShow) {
      for (let i = 1; i <= totalPages; i++) {
        pages.push(i);
      }
    } else {
      if (currentPage <= 3) {
        for (let i = 1; i <= 4; i++) {
          pages.push(i);
        }
        pages.push('...');
        pages.push(totalPages);
      } else if (currentPage >= totalPages - 2) {
        pages.push(1);
        pages.push('...');
        for (let i = totalPages - 3; i <= totalPages; i++) {
          pages.push(i);
        }
      } else {
        pages.push(1);
        pages.push('...');
        pages.push(currentPage - 1);
        pages.push(currentPage);
        pages.push(currentPage + 1);
        pages.push('...');
        pages.push(totalPages);
      }
    }

    return pages;
  };

  const onChangePage = (page: number | string) => {
    if (
      typeof page === 'number' &&
      page >= 1 &&
      page <= totalPages &&
      page !== currentPage
    ) {
      goToPage(page);
    }
  };

  return (
    <Fragment>
      {totalPages > 1 && (
        <div className="mt-6 flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="text-sm text-slate-600 dark:text-slate-400">
            Page {currentPage} of {totalPages}
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={() => onChangePage(currentPage - 1)}
              disabled={currentPage === 1}
              className={`flex items-center gap-1 px-3 py-2 rounded-lg text-sm font-medium transition-all ${
                currentPage === 1
                  ? 'bg-slate-200 dark:bg-slate-700 text-slate-400 dark:text-slate-500 cursor-not-allowed'
                  : 'bg-blue-500 dark:bg-blue-600 text-white hover:bg-blue-600 dark:hover:bg-blue-500 hover:scale-105'
              }`}
            >
              <ChevronLeft className="w-4 h-4" />
              <span className="hidden sm:inline">Previous</span>
            </button>

            <div className="flex items-center gap-1">
              {getPageNumbers().map((pageNum, idx) => (
                <Fragment key={idx}>
                  {pageNum === '...' ? (
                    <span className="px-3 py-2 text-slate-500 dark:text-slate-400">
                      ...
                    </span>
                  ) : (
                    <button
                      onClick={() => onChangePage(pageNum)}
                      className={`px-4 py-2 rounded-lg text-sm font-medium transition-all hover:scale-105 ${
                        currentPage === pageNum
                          ? 'bg-blue-500 dark:bg-blue-600 text-white'
                          : 'bg-slate-200 dark:bg-slate-700 text-slate-700 dark:text-slate-300 hover:bg-slate-300 dark:hover:bg-slate-600'
                      }`}
                    >
                      {pageNum}
                    </button>
                  )}
                </Fragment>
              ))}
            </div>

            <button
              onClick={() => onChangePage(currentPage + 1)}
              disabled={currentPage === totalPages}
              className={`flex items-center gap-1 px-3 py-2 rounded-lg text-sm font-medium transition-all ${
                currentPage === totalPages
                  ? 'bg-slate-200 dark:bg-slate-700 text-slate-400 dark:text-slate-500 cursor-not-allowed'
                  : 'bg-blue-500 dark:bg-blue-600 text-white hover:bg-blue-600 dark:hover:bg-blue-500 hover:scale-105'
              }`}
            >
              <span className="hidden sm:inline">Next</span>
              <ChevronRight className="w-4 h-4" />
            </button>
          </div>
        </div>
      )}
    </Fragment>
  );
};

export default TransactionPagination;
