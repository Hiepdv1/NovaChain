'use client';

import React, { memo } from 'react';

const NetworkOverviewSkeleton = ({}) => {
  const cards = Array(5).fill(null);

  return (
    <div className="mb-8 animate-pulse">
      <div className="h-7 w-48 bg-gray-300 dark:bg-gray-700 rounded-lg mb-6" />

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6">
        {cards.map((_, idx) => (
          <div
            key={idx}
            className="glass-card bg-gradient-glass dark:bg-gradient-glass-dark rounded-2xl p-6 border border-white/20 dark:border-gray-700/50"
          >
            <div className="flex items-center justify-between mb-4">
              <div className="h-4 w-24 bg-gray-300 dark:bg-gray-700 rounded" />
              <div className="w-8 h-8 bg-gray-200 dark:bg-gray-800 rounded-lg" />
            </div>
            <div className="h-6 w-24 bg-gray-300 dark:bg-gray-700 rounded mb-2" />
            <div className="h-3 w-16 bg-gray-200 dark:bg-gray-800 rounded" />
          </div>
        ))}
      </div>

      <div className="mt-6 h-[2px] w-full bg-gradient-to-r from-transparent via-gray-300 dark:via-gray-700 to-transparent animate-[pulse_2s_infinite]" />
    </div>
  );
};

export default memo(NetworkOverviewSkeleton);
