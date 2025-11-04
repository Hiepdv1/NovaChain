import { AlertCircle, CheckCircle2, XCircle } from 'lucide-react';

const StatusBadge = ({
  status,
}: {
  status: 'success' | 'pending' | 'failed';
}) => {
  const config = {
    success: {
      icon: CheckCircle2,
      className: 'bg-gradient-to-r from-emerald-500 to-green-500 text-white',
      label: 'Success',
    },
    pending: {
      icon: AlertCircle,
      className: 'bg-gradient-to-r from-amber-500 to-orange-500 text-white',
      label: 'Pending',
    },
    failed: {
      icon: XCircle,
      className: 'bg-gradient-to-r from-red-500 to-pink-500 text-white',
      label: 'Failed',
    },
  };

  const { icon: Icon, className, label } = config[status];

  return (
    <span
      className={`inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-semibold ${className} shadow-sm`}
    >
      <Icon className="w-3.5 h-3.5" />
      {label}
    </span>
  );
};

export default StatusBadge;
