import { CheckCircle2, ExternalLink } from 'lucide-react';
import { ComponentType } from 'react';
import CopyButton from './tx-detail-btn-copy';
import Link from 'next/link';

interface InfoRowProps {
  icon: ComponentType<{ className: string }>;
  label: string;
  value: string;
  copyable?: string;
  link?: string;
  status?: string;
}

const InfoRow = ({
  icon: Icon,
  label,
  value,
  copyable,
  link,
  status,
}: InfoRowProps) => {
  return (
    <div className="flex flex-col sm:flex-row sm:items-center py-4 border-b border-gray-100 dark:border-gray-700 last:border-0">
      <div className="flex items-center gap-2 sm:w-1/3 mb-2 sm:mb-0">
        {Icon && (
          <div className="p-1.5 rounded-lg bg-gray-100 dark:bg-gray-700">
            <Icon className="w-4 h-4 text-gray-600 dark:text-gray-400" />
          </div>
        )}
        <span className="text-sm font-medium text-gray-600 dark:text-gray-400">
          {label}
        </span>
      </div>
      <div className="sm:w-2/3 flex items-center gap-2 flex-wrap">
        {status ? (
          <span
            className={`inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-sm font-semibold ${
              status === 'success'
                ? 'bg-green-100 dark:bg-green-900/30 text-green-700 dark:text-green-400'
                : status === 'pending'
                ? 'bg-yellow-100 dark:bg-yellow-900/30 text-yellow-700 dark:text-yellow-400'
                : 'bg-red-100 dark:bg-red-900/30 text-red-700 dark:text-red-400'
            }`}
          >
            <CheckCircle2 className="w-4 h-4" />
            {value}
          </span>
        ) : link ? (
          <Link
            href={link}
            className="text-sm font-mono text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 transition-colors inline-flex items-center gap-1"
          >
            {value}
            <ExternalLink className="w-3.5 h-3.5" />
          </Link>
        ) : (
          <span className="text-sm font-mono text-gray-900 dark:text-gray-100 break-all">
            {value}
          </span>
        )}
        {copyable && <CopyButton text={copyable} label={label} />}
      </div>
    </div>
  );
};

export default InfoRow;
