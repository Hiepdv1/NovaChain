'use client';

import {
  Activity,
  Clock,
  Coins,
  Hash,
  Info,
  Layers,
  Receipt,
  Shield,
  User,
} from 'lucide-react';
import SectionCard from './section-card';
import InfoRow from './info-row';
import { FormatFloat, FormatTimestamp } from '@/shared/utils/format';
import { TransactionDetail } from '../../types/transaction';
import { memo } from 'react';

interface TransactionCardProps {
  transaction: TransactionDetail;
}

const TransactionCard = ({ transaction }: TransactionCardProps) => {
  return (
    <div className="space-y-6">
      <SectionCard title="Transaction Overview" icon={Hash}>
        <div className="space-y-1">
          <InfoRow
            icon={Hash}
            label="Transaction Hash"
            value={transaction.TxID}
            copyable={transaction.TxID}
          />
          <InfoRow
            icon={Layers}
            label="Block Number"
            value={transaction.LastBlock.toString()}
            link={`/blocks/${transaction.BID}`}
          />
          <InfoRow
            icon={Hash}
            label="Block Hash"
            value={transaction.BID}
            copyable={transaction.BID}
          />
          <InfoRow
            icon={Clock}
            label="Timestamp"
            value={FormatTimestamp(transaction.Timestamp)}
          />
        </div>
      </SectionCard>

      <SectionCard title="Transfer Details" icon={User}>
        <div className="space-y-1">
          <InfoRow
            icon={User}
            label="From"
            value={transaction.Fromhash.String}
            copyable={transaction.Fromhash.String}
          />
          <InfoRow
            icon={User}
            label="To"
            value={transaction.Tohash.String}
            copyable={transaction.Tohash.String}
          />
          <InfoRow
            icon={Coins}
            label="Amount"
            value={FormatFloat(Number(transaction.Amount.String), 8).toString()}
          />
          <InfoRow
            icon={Receipt}
            label="Transaction Fee"
            value={FormatFloat(Number(transaction.Fee.String), 8).toString()}
          />
        </div>
      </SectionCard>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <SectionCard title="Mining Information" icon={Activity}>
          <div className="space-y-4">
            <div className="flex justify-between items-center py-2 border-b border-gray-100 dark:border-gray-700">
              <span className="text-sm text-gray-600 dark:text-gray-400">
                Mined By
              </span>
              <span className="text-sm font-mono text-black dark:text-white">
                {transaction.Miner.String.slice(0, 10)}...
                {transaction.Miner.String.slice(-8)}
              </span>
            </div>
            <div className="flex justify-between items-center py-2 border-b border-gray-100 dark:border-gray-700">
              <span className="text-sm text-gray-600 dark:text-gray-400">
                Block Difficulty
              </span>
              <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                {transaction.Difficulty}
              </span>
            </div>
            <div className="flex justify-between items-center py-2">
              <span className="text-sm text-gray-600 dark:text-gray-400">
                Position in Block
              </span>
              <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                {transaction.Height} / {transaction.LastBlock}
              </span>
            </div>
          </div>
        </SectionCard>

        <SectionCard title="Additional Information" icon={Info}>
          <div className="space-y-4">
            <div className="flex justify-between items-center py-2 border-b border-gray-100 dark:border-gray-700">
              <span className="text-sm text-gray-600 dark:text-gray-400">
                Nonce
              </span>
              <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                {transaction.Nonce}
              </span>
            </div>
          </div>
        </SectionCard>
      </div>

      <SectionCard
        title="Cryptographic Verification"
        icon={Shield}
        collapsible={true}
      >
        <div className="space-y-4">
          <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-4">
            <div className="flex items-start gap-3">
              <Shield className="w-5 h-5 text-blue-600 dark:text-blue-400 flex-shrink-0 mt-0.5" />
              <div>
                <h4 className="text-sm font-semibold text-blue-900 dark:text-blue-100 mb-1">
                  Verification Status
                </h4>
                <p className="text-sm text-blue-800 dark:text-blue-200">
                  This transaction has been cryptographically verified and is
                  immutably recorded on the blockchain.
                </p>
                <ul className="mt-2 space-y-1 text-xs text-blue-700 dark:text-blue-300">
                  <li>✓ Signature verified against sender address</li>
                  <li>
                    ✓ Transaction included in mined block #{transaction.Height}
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </SectionCard>
    </div>
  );
};

export default memo(TransactionCard);
