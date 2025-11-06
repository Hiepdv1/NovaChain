import CopyButton from '@/components/button/copyButton';
import { Miner } from '../types/miner';
import { FormatTimestamp, TruncateHash } from '@/shared/utils/format';
import RankBadge from './miner-rank';
import { memo } from 'react';
import { Award, Clock, TrendingUp } from 'lucide-react';

interface MinerCardProps {
  miner: Miner;
  rank: number;
}

const MinerCard = ({ miner, rank }: MinerCardProps) => {
  return (
    <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-sm border border-gray-200 dark:border-gray-700 p-6 hover:shadow-xl hover:scale-[1.01] transition-all duration-300">
      <div className="flex items-center gap-4 mb-6">
        <RankBadge rank={rank} />

        <div className="flex-1 min-w-0">
          <div className="flex items-center gap-2 mb-1">
            <span className="text-lg font-mono font-bold text-gray-900 dark:text-gray-100 truncate">
              {TruncateHash(miner.MinerPubkey)}
            </span>
            <CopyButton type="icon" text={miner.MinerPubkey} />
          </div>

          {rank <= 3 && (
            <div className="flex items-center gap-1 text-sm">
              <Award
                className={`w-4 h-4 ${
                  rank === 1
                    ? 'text-yellow-500'
                    : rank === 2
                    ? 'text-gray-400'
                    : 'text-orange-500'
                }`}
              />
              <span
                className={`font-semibold ${
                  rank === 1
                    ? 'text-yellow-600 dark:text-yellow-400'
                    : rank === 2
                    ? 'text-gray-600 dark:text-gray-400'
                    : 'text-orange-600 dark:text-orange-400'
                }`}
              >
                Top {rank} Miner
              </span>
            </div>
          )}
        </div>

        <div className="text-right">
          <div className="text-2xl font-bold text-purple-600 dark:text-purple-400">
            {miner.MinedBlocks}
          </div>
          <div className="text-xs text-gray-500 dark:text-gray-400">blocks</div>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-2 gap-3">
        <div className="bg-gradient-to-br from-purple-50 to-indigo-50 dark:from-purple-900/20 dark:to-indigo-900/20 rounded-xl p-3 border border-purple-100 dark:border-purple-800">
          <div className="flex items-center gap-2 mb-1">
            <TrendingUp className="w-4 h-4 text-purple-600 dark:text-purple-400" />
            <span className="text-xs font-medium text-purple-700 dark:text-purple-300">
              Network Share
            </span>
          </div>
          <div className="text-lg font-bold text-purple-900 dark:text-purple-100">
            {miner.NetworkSharePercent}%
          </div>
        </div>

        <div className="bg-gradient-to-br from-blue-50 to-cyan-50 dark:from-blue-900/20 dark:to-cyan-900/20 rounded-xl p-3 border border-blue-100 dark:border-blue-800">
          <div className="flex items-center gap-2 mb-1">
            <Clock className="w-4 h-4 text-blue-600 dark:text-blue-400" />
            <span className="text-xs font-medium text-blue-700 dark:text-blue-300">
              Last Block
            </span>
          </div>
          <div className="text-lg font-bold text-blue-900 dark:text-blue-100">
            {FormatTimestamp(miner.LastMinedAt)}
          </div>
        </div>
      </div>

      {/* Progress Bar */}
      <div className="mt-4">
        <div className="flex items-center justify-between mb-2">
          <span className="text-xs font-medium text-gray-600 dark:text-gray-400">
            Mining Activity
          </span>
          <span className="text-xs font-bold text-gray-900 dark:text-gray-100">
            {miner.MinedBlocks} / {miner.TotalBlocksNetwork}
          </span>
        </div>
        <div className="w-full h-2 bg-gray-200 dark:bg-gray-700 rounded-full overflow-hidden">
          <div
            className="h-full bg-gradient-to-r from-purple-500 to-indigo-500 rounded-full transition-all duration-500"
            style={{
              width: `${(miner.MinedBlocks * 100) / miner.TotalBlocksNetwork}%`,
            }}
          />
        </div>
      </div>
    </div>
  );
};

export default memo(MinerCard);
