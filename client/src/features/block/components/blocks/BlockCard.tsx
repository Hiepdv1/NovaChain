import {
  FormatFloat,
  FormatSize,
  FormatTimestamp,
  TruncateHash,
} from '@/shared/utils/format';
import { BlockItem } from '../../types/block';
import { memo, useState } from 'react';
import { Check, Copy } from 'lucide-react';
import Link from 'next/link';
import { useHistoryContext } from '@/components/providers/History-provider';

interface BlockCardProps {
  block: BlockItem;
}

const BlockCard = ({ block }: BlockCardProps) => {
  const [copiedHash, setCopiedHash] = useState<string | null>(null);
  const { push } = useHistoryContext();

  const onCopy = (hash: string) => {
    navigator.clipboard.writeText(hash);
    setCopiedHash(hash);
    setTimeout(() => setCopiedHash(null), 2000);
  };

  return (
    <Link
      href={`/blocks/${block.BID}`}
      onClick={() =>
        push({
          path: '/blocks',
          title: 'Back',
        })
      }
      className="block rounded-xl p-6 transition-all duration-300 hover:scale-105 hover:shadow-xl bg-white dark:bg-gray-800 shadow-md dark:shadow-lg animate-block-fade-in"
    >
      <div className="flex items-center justify-between mb-4">
        <a
          href="#"
          className="text-2xl font-bold text-blue-500 hover:text-blue-400 transition-colors"
        >
          #{block.Height}
        </a>
        <span className="text-xs px-2 py-1 rounded-full bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300">
          {FormatTimestamp(block.Timestamp)}
        </span>
      </div>

      <div className="mb-4">
        <label className="text-xs font-semibold uppercase text-gray-500 dark:text-gray-500">
          Block Hash
        </label>
        <div className="flex items-center gap-2 mt-1">
          <code className="text-sm font-mono flex-1 text-gray-700 dark:text-gray-300">
            {TruncateHash(block.BID)}
          </code>
          <button
            onClick={() => onCopy(block.BID)}
            className={`${
              copiedHash ? '' : 'cursor-pointer'
            } p-1.5 rounded transition-colors hover:bg-gray-100 dark:hover:bg-gray-700`}
          >
            {copiedHash === block.BID ? (
              <Check size={16} className="text-green-500" />
            ) : (
              <Copy size={16} />
            )}
          </button>
        </div>
      </div>

      <div className="space-y-3">
        <div className="flex justify-between items-center">
          <span className="text-gray-600 dark:text-gray-400">Miner</span>
          <span className="text-blue-500 hover:text-blue-400 font-mono text-sm transition-colors">
            {TruncateHash(block.PubKeyHash, 6, 4)}
          </span>
        </div>

        <div className="flex justify-between items-center">
          <span className="text-gray-600 dark:text-gray-400">Transactions</span>
          <span className="font-semibold">{block.TxCount}</span>
        </div>

        <div className="flex justify-between items-center">
          <span className="text-gray-600 dark:text-gray-400">Size</span>
          <span className="font-semibold">{FormatSize(block.Size)}</span>
        </div>

        <div className="flex justify-between items-center pt-3 border-t border-gray-200 dark:border-gray-700">
          <span className="text-gray-600 dark:text-gray-400">Reward</span>
          <span className="font-bold text-green-500">
            {FormatFloat(Number(block.Value), 2)} CCC
          </span>
        </div>
      </div>
    </Link>
  );
};

export default memo(BlockCard);
