import { Check, Copy, ExternalLink, Hash, Search, X, Zap } from 'lucide-react';
import { BlockDetail } from '../types/block';
import {
  ChangeEvent,
  Fragment,
  memo,
  useCallback,
  useRef,
  useState,
} from 'react';
import { useSearchTransactionsByBlockQuery } from '../hook/useBlockQuery';
import { CACHE_TIME } from '@/shared/constants/ttl';
import { FormatTimestamp, TruncateHash } from '@/shared/utils/format';
import ErrorState from '@/components/errorState';
import TransactionListSkeleton from './tx-list-Loading';
import TransactionPagination from './tx-pagination';

interface TransactionListProps {
  blockData: BlockDetail;
}

const TransactionList = ({ blockData }: TransactionListProps) => {
  const [searchQuery, setSearchQuery] = useState({
    value: '',
    isSearch: false,
  });
  const [copiedHash, setCopiedHash] = useState<string | null>(null);
  const delay = useRef<NodeJS.Timeout | null>(null);
  const [page, setPage] = useState(1);
  const totalPages = Math.ceil(
    blockData.Transactions.Meta.total / blockData.Transactions.Meta.limit,
  );

  const { isError, isLoading, isFetching, data, error, refetch } =
    useSearchTransactionsByBlockQuery(
      {
        b_hash: blockData.BID,
        q: searchQuery.value,
        page,
        limit: 10,
      },
      {
        enabled:
          !!blockData.BID &&
          blockData.TxCount > 0 &&
          searchQuery.value.length <= 64 &&
          searchQuery.isSearch,
        staleTime: 0,
        gcTime: CACHE_TIME,
        refetchOnWindowFocus: false,
        retry: false,
      },
    );

  const copyToClipboard = (text: string, id: string) => {
    navigator.clipboard.writeText(text);
    setCopiedHash(id);
    setTimeout(() => setCopiedHash(null), 2000);
  };

  const onChangeInput = (e: ChangeEvent<HTMLInputElement>) => {
    if (delay.current) {
      clearTimeout(delay.current);
    }
    setSearchQuery({ value: e.target.value.trim(), isSearch: false });

    delay.current = setTimeout(() => {
      setSearchQuery((prev) => ({ ...prev, isSearch: true }));
    }, 500);
  };

  const onGoToPage = useCallback((page: number) => {
    setPage(page);
    setSearchQuery({
      value: '',
      isSearch: true,
    });
  }, []);

  const onRetry = useCallback(() => {
    refetch();
  }, [refetch]);

  const clearSearch = () => {
    setSearchQuery({
      value: '',
      isSearch: false,
    });
  };

  const filteredTxs = data?.data || blockData.Transactions.Data;

  return (
    <div className="rounded-xl p-6 transition-all bg-gradient-to-br from-white to-purple-50 dark:from-slate-800 dark:to-slate-700 shadow-lg shadow-purple-200/50 dark:shadow-blue-500/10">
      <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
        <h2 className="text-2xl font-bold flex items-center gap-2 text-slate-900 dark:text-slate-50">
          <Zap className="w-6 h-6 text-blue-500 dark:text-blue-400" />
          Transactions ({blockData.TxCount})
        </h2>

        {/* Search Bar */}
        <div className="relative w-72 max-sm:w-full lg:w-96">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 w-4 h-4 text-slate-400 dark:text-slate-500" />
          <input
            type="text"
            placeholder="Search transactions..."
            value={searchQuery.value}
            onChange={onChangeInput}
            className="w-full pl-9 pr-9 py-3 text-sm rounded-lg border transition-all focus:outline-none focus:ring-2 bg-white dark:bg-slate-900/50 border-blue-200 dark:border-slate-600 text-slate-900 dark:text-slate-100 placeholder:text-slate-400 dark:placeholder:text-slate-500 focus:border-blue-400 dark:focus:border-blue-400 focus:ring-blue-300/30 dark:focus:ring-blue-500/30"
          />
          {searchQuery.value && (
            <button
              onClick={clearSearch}
              className="cursor-pointer absolute right-1.5 top-1/2 transform -translate-y-1/2 p-1 rounded-full transition-all hover:scale-110 hover:bg-slate-200 dark:hover:bg-slate-700"
            >
              <X className="w-3.5 h-3.5 text-slate-600 dark:text-slate-400" />
            </button>
          )}
        </div>
      </div>

      <div className="my-4 text-sm text-slate-600 dark:text-slate-400">
        {searchQuery.value && (
          <span>
            Found {data?.data.length} transaction
            {data?.data.length !== 1 ? 's' : ''} matching &quot;
            {searchQuery.value}
            &quot;
          </span>
        )}
      </div>

      {!isError && !isFetching && !isLoading && filteredTxs.length > 0 ? (
        <div className="space-y-4">
          {filteredTxs.map((tx, idx) => {
            return (
              <div
                key={idx}
                className="rounded-lg p-4 transition-all hover:scale-[1.01] border bg-gradient-to-br from-white to-blue-50 dark:from-slate-900/50 dark:to-slate-800  border-blue-200 dark:border-slate-600 hover:border-blue-400 dark:hover:border-blue-400/50 hover:shadow-lg shadow-blue-300/30 dark:shadow-blue-500/10"
              >
                <div className="flex items-start justify-between mb-3">
                  <div className="flex items-center gap-2 flex-1 min-w-0">
                    <Hash className="w-4 h-4 text-blue-500 dark:text-blue-400 flex-shrink-0" />
                    <code className="text-sm font-mono truncate text-blue-500 dark:text-blue-400 hover:text-blue-600 dark:hover:text-blue-300 cursor-pointer">
                      {TruncateHash(tx.TxID, 12, 10)}
                    </code>
                    <button
                      onClick={() => copyToClipboard(tx.TxID, `tx-${idx}`)}
                      className="p-1 rounded transition-all hover:scale-110 hover:bg-blue-100 dark:hover:bg-slate-700"
                    >
                      {copiedHash === `tx-${idx}` ? (
                        <Check className="w-4 h-4 text-green-500 dark:text-green-400" />
                      ) : (
                        <Copy className="w-4 h-4 text-slate-600 dark:text-slate-400" />
                      )}
                    </button>
                  </div>
                  <span className="text-xs text-slate-600 dark:text-slate-400">
                    {FormatTimestamp(tx.CreateAt)}
                  </span>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 text-sm">
                  <div>
                    <span className="block mb-1 text-slate-600 dark:text-slate-400">
                      From
                    </span>
                    <code className="font-mono truncate text-xs text-slate-700 dark:text-slate-300">
                      {typeof tx.Fromhash === 'object' && tx.Fromhash !== null
                        ? tx.Fromhash.String
                        : tx.Fromhash || 'N/A'}
                    </code>
                  </div>
                  <div>
                    <span className="block mb-1 text-slate-600 dark:text-slate-400">
                      To
                    </span>
                    <code className="font-mono truncate text-xs text-slate-700 dark:text-slate-300">
                      {typeof tx.Tohash === 'object' && tx.Tohash !== null
                        ? tx.Tohash.String
                        : tx.Tohash || 'N/A'}
                    </code>
                  </div>
                </div>

                <div className="flex items-center justify-between mt-3 pt-3 border-t border-blue-200 dark:border-slate-600">
                  <div className="flex items-center gap-4">
                    <div>
                      <span className="text-xs text-slate-600 dark:text-slate-400">
                        Value
                      </span>
                      <p className="font-semibold truncate text-green-600 dark:text-green-400">
                        {typeof tx.Amount === 'object' && tx.Amount !== null
                          ? tx.Amount.String
                          : tx.Amount || '0'}
                      </p>
                    </div>
                    <div>
                      <span className="text-xs text-slate-600 dark:text-slate-400">
                        Fee
                      </span>
                      <p className="font-semibold truncate text-slate-700 dark:text-slate-300">
                        {typeof tx.Fee === 'object' && tx.Fee !== null
                          ? tx.Fee.String
                          : tx.Fee || '0'}
                      </p>
                    </div>
                  </div>
                  <button className="cursor-pointer flex items-center gap-1 px-3 py-2 rounded-lg text-xs transition-all hover:scale-105 bg-gradient-primary dark:bg-blue-500/20 text-white">
                    View Details
                    <ExternalLink className="w-3 h-3" />
                  </button>
                </div>
              </div>
            );
          })}
          <TransactionPagination
            currentPage={data?.meta.currentPage || 1}
            totalPages={
              (data && Math.ceil(data.meta.total / data.meta.limit)) ||
              totalPages
            }
            goToPage={onGoToPage}
          />
        </div>
      ) : (
        <Fragment>
          {(isLoading || isFetching) && <TransactionListSkeleton />}
          {isError && error.statusCode === 404 && (
            <div className="text-center py-12 text-slate-600 dark:text-slate-400">
              <Search className="w-12 h-12 mx-auto mb-3 opacity-50" />
              <p className="text-lg font-medium">No transactions found</p>
              <p className="text-sm mt-1">Try adjusting your search query</p>
            </div>
          )}
          {isError && error.statusCode !== 404 && (
            <ErrorState onRetry={onRetry} />
          )}
        </Fragment>
      )}
    </div>
  );
};

export default memo(TransactionList);
