'use client';

import { SearchX } from 'lucide-react';
import Link from 'next/link';

interface NotFoundStateProps {
  title?: string;
  message?: string;
}

const NotFoundState = ({
  title = 'Not Found',
  message = 'Sorry, we couldn’t find what you were looking for. It may have been moved, deleted, or never existed.',
}: NotFoundStateProps) => {
  return (
    <div className="flex flex-col items-center justify-center py-20 px-6 text-center">
      {/* Icon + Glow background */}
      <div className="relative mb-8">
        <div className="absolute inset-0 bg-gradient-to-r from-indigo-500 to-purple-500 rounded-full blur-3xl opacity-25 animate-pulse"></div>
        <div className="relative p-8 bg-gradient-to-br from-gray-100 to-gray-200 dark:from-gray-800 dark:to-gray-700 rounded-full shadow-lg">
          <SearchX className="w-16 h-16 text-gray-400 dark:text-gray-500" />
        </div>
      </div>

      <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-3">
        {title}
      </h2>

      <p className="text-gray-600 dark:text-gray-400 max-w-md leading-relaxed mb-8">
        {message}
      </p>

      <Link
        href="/"
        className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-primary text-white font-medium shadow hover:opacity-90 transition"
      >
        Go Home
      </Link>
    </div>
  );
};

export default NotFoundState;
