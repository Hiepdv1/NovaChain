import { Loader2, XCircle } from 'lucide-react';
import { memo } from 'react';

interface StatusBadgeProps {
  status: 'pending' | 'failed';
  fee: number;
}

const TxPendingStatusBadge = ({ status, fee }: StatusBadgeProps) => {
  const config = {
    pending: {
      icon: Loader2,
      className: 'bg-gradient-to-r from-amber-500 to-orange-500 text-white',
      label: 'Pending',
    },
    failed: {
      icon: XCircle,
      className: 'bg-gradient-to-r from-red-500 to-pink-500 text-white',
      label: 'Failed',
    },
  };

  const { icon: Icon, className, label } = config[status] || config.pending;

  return (
    <div className="flex flex-col items-end gap-2">
      <span
        className={`inline-flex items-center gap-2 px-4 py-2 rounded-full text-sm font-bold ${className} shadow-lg`}
      >
        <Icon className="w-4 h-4 animate-spin" />
        {label}
      </span>
      {fee && (
        <div className="px-3 py-1 bg-gray-100 dark:bg-gray-700 rounded-full text-xs font-semibold text-gray-700 dark:text-gray-300">
          Fee: {fee}
        </div>
      )}
    </div>
  );
};

export default memo(TxPendingStatusBadge);
