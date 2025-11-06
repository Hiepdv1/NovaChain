import { memo, ComponentType } from 'react';

export interface StatCardProps {
  Icon: ComponentType<{ className?: string }>;
  label: string;
  value: number;
  subValue?: string;
  color: 'blue' | 'green' | 'purple' | 'orange';
}

const StatCard = ({ Icon, label, value, subValue, color }: StatCardProps) => {
  const colorClasses = {
    blue: 'bg-gradient-to-br from-blue-500 to-blue-600',
    green: 'bg-gradient-to-br from-green-500 to-green-600',
    purple: 'bg-gradient-to-br from-purple-500 to-purple-600',
    orange: 'bg-gradient-to-br from-orange-500 to-orange-600',
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl p-5 border border-gray-200 dark:border-gray-700 hover:shadow-md transition-all">
      <div
        className={`inline-flex p-2.5 rounded-lg ${colorClasses[color]} mb-3`}
      >
        <Icon className="w-5 h-5 text-white" />
      </div>
      <div className="text-sm text-gray-600 dark:text-gray-400 mb-1">
        {label}
      </div>
      <div className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-1">
        {value}
      </div>
      {subValue && (
        <div className="text-sm text-gray-500 dark:text-gray-400">
          {subValue}
        </div>
      )}
    </div>
  );
};

export default memo(StatCard);
