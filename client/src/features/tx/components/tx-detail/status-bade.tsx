import { AlertCircle, CheckCircle2, Clock } from 'lucide-react';

type StatusType = 'confirmed' | 'pending' | 'failed';

const StatusBadge = ({ status }: { status: StatusType }) => {
  const config = {
    confirmed: {
      icon: CheckCircle2,
      label: 'Confirmed',
      className: 'bg-gradient-to-r from-green-500 to-emerald-500 text-white',
    },
    pending: {
      icon: Clock,
      label: 'Pending',
      className: 'bg-gradient-to-r from-yellow-500 to-orange-500 text-white',
    },
    failed: {
      icon: AlertCircle,
      label: 'Failed',
      className: 'bg-gradient-to-r from-red-500 to-pink-500 text-white',
    },
  };

  const { icon: Icon, label, className } = config[status] || config.failed;

  return (
    <div
      className={`inline-flex items-center gap-2 px-4 py-2 rounded-full ${className} shadow-lg`}
    >
      <Icon className="w-5 h-5" />
      <span className="font-semibold">{label}</span>
    </div>
  );
};

export default StatusBadge;
