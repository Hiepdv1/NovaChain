'use client';

import { ChevronDown, ChevronUp } from 'lucide-react';
import { ComponentType, useState } from 'react';

interface SectionCardProps {
  title: string;
  icon: ComponentType<{ className: string }>;
  children: React.ReactNode;
  collapsible?: boolean;
}

const SectionCard = ({
  title,
  icon: Icon,
  children,
  collapsible = false,
}: SectionCardProps) => {
  const [isExpanded, setIsExpanded] = useState(true);

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden transition-colors">
      <div
        className={`px-6 py-4 border-b border-gray-200 dark:border-gray-700 ${
          collapsible
            ? 'cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700/50'
            : ''
        }`}
        onClick={() => collapsible && setIsExpanded(!isExpanded)}
      >
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            {Icon && (
              <div className="p-2 rounded-lg bg-gradient-to-br from-blue-500 to-purple-600">
                <Icon className="w-5 h-5 text-white" />
              </div>
            )}
            <h2 className="text-lg font-bold text-gray-900 dark:text-gray-100">
              {title}
            </h2>
          </div>
          {collapsible && (
            <button className="p-1 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition-colors">
              {isExpanded ? (
                <ChevronUp className="w-5 h-5 text-gray-600 dark:text-gray-400" />
              ) : (
                <ChevronDown className="w-5 h-5 text-gray-600 dark:text-gray-400" />
              )}
            </button>
          )}
        </div>
      </div>
      {isExpanded && <div className="p-6">{children}</div>}
    </div>
  );
};

export default SectionCard;
