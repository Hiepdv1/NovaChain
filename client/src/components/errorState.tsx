'use client';

import { AlertTriangle, RefreshCw } from 'lucide-react';
import { memo, useState } from 'react';
import { useRouter } from 'next/navigation';

interface ErrorStateProps {
  title?: string;
  message?: string;
  onRetry?: () => Promise<void> | void;
  retryText?: string;
  icon?: React.ElementType;
  gradient?: string;
}

const ErrorState = ({
  title = 'Something went wrong',
  message = 'An unexpected error occurred. Please try again later.',
  onRetry,
  retryText = 'Try again',
  icon: Icon = AlertTriangle,
  gradient = 'from-red-500 to-orange-500',
}: ErrorStateProps) => {
  const [isRetrying, setIsRetrying] = useState(false);
  const router = useRouter();

  const handleRetry = async () => {
    setIsRetrying(true);
    try {
      if (onRetry) {
        await onRetry();
      } else {
        router.refresh();
      }
    } finally {
      setTimeout(() => setIsRetrying(false), 400);
    }
  };

  return (
    <div className="flex flex-col items-center justify-center py-20 px-4">
      <div className="relative">
        <div
          className={`absolute inset-0 bg-gradient-to-r ${gradient} rounded-full blur-3xl opacity-20 animate-pulse`}
        />
        <div className="relative p-8 bg-gradient-to-br from-red-100 to-orange-100 dark:from-gray-800 dark:to-gray-700 rounded-full mb-6">
          <Icon className="w-16 h-16 text-red-500 dark:text-red-400" />
        </div>
      </div>
      <h3 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-3">
        {title}
      </h3>
      <p className="text-gray-600 dark:text-gray-400 text-center max-w-md leading-relaxed mb-6">
        {message}
      </p>
      <button
        onClick={handleRetry}
        disabled={isRetrying}
        className="flex items-center gap-2 px-6 py-2 rounded-xl bg-blue-600 hover:bg-blue-700 text-white font-medium transition disabled:opacity-50"
      >
        <RefreshCw className={`w-4 h-4 ${isRetrying ? 'animate-spin' : ''}`} />
        {isRetrying ? 'Retrying...' : retryText}
      </button>
    </div>
  );
};

export default memo(ErrorState);
