'use client';

import { useEffect, useState } from 'react';

const RootLoader = () => {
  const [messageIndex, setMessageIndex] = useState(0);
  const [dotCount, setDotCount] = useState(0);

  const messages = [
    'Connecting to network',
    'Loading blockchain data',
    'Fetching transactions',
    'Preparing interface',
  ];

  useEffect(() => {
    const messageInterval = setInterval(() => {
      setMessageIndex((prev) => (prev + 1) % messages.length);
    }, 1000);

    const dotInterval = setInterval(() => {
      setDotCount((prev) => (prev + 1) % 4);
    }, 800);

    return () => {
      clearInterval(messageInterval);
      clearInterval(dotInterval);
    };
  }, [messages.length]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-white dark:bg-gray-900 transition-colors duration-500">
      <div className="absolute inset-0 opacity-30">
        <div
          className="absolute top-1/4 left-1/4 w-64 h-64 bg-blue-500/15 dark:bg-blue-500/10 rounded-full blur-3xl animate-pulse"
          style={{ animationDuration: '4s' }}
        ></div>
        <div
          className="absolute bottom-1/4 right-1/4 w-48 h-48 bg-purple-500/15 dark:bg-purple-500/10 rounded-full blur-3xl animate-pulse"
          style={{ animationDuration: '3s', animationDelay: '1s' }}
        ></div>
      </div>

      <div className="relative z-10 text-center">
        <div className="mb-8 flex justify-center">
          <div className="relative w-20 h-20 bg-gradient-to-br from-blue-500 to-purple-500 dark:from-blue-600 dark:to-purple-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-500/20 dark:shadow-blue-500/25 transition-all duration-500">
            <div
              className="absolute -inset-1 bg-gradient-to-r from-blue-400/60 to-purple-400/60 dark:from-blue-500/50 dark:to-purple-500/50 rounded-2xl opacity-75 animate-spin"
              style={{ animationDuration: '3s' }}
            ></div>

            <svg
              className="relative z-10 w-10 h-10 text-white animate-pulse"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.94-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z" />
            </svg>
          </div>
        </div>

        <h1 className="text-3xl font-bold mb-2 bg-gradient-to-r from-blue-600 to-purple-600 dark:from-blue-400 dark:to-purple-400 bg-clip-text text-transparent transition-colors duration-500">
          CryptoChain
        </h1>

        <p className="text-sm mb-12 text-gray-600 dark:text-gray-400 transition-colors duration-500">
          Blockchain Explorer
        </p>

        <div className="mb-8 flex justify-center">
          <div className="flex items-center space-x-3">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="flex items-center">
                <div
                  className="w-3 h-3 bg-blue-600 dark:bg-blue-500 rounded-full transition-colors duration-500"
                  style={{
                    animation: `pulse 2s ease-in-out infinite`,
                    animationDelay: `${i * 0.4}s`,
                    opacity: 0.4,
                  }}
                />
                {i < 4 && (
                  <div
                    className="w-8 h-0.5 mx-1 bg-gradient-to-r from-blue-400/60 to-purple-400/60 dark:from-blue-500/50 dark:to-purple-500/50 transition-colors duration-500"
                    style={{
                      animation: `pulse 2s ease-in-out infinite`,
                      animationDelay: `${i * 0.4 + 0.2}s`,
                    }}
                  />
                )}
              </div>
            ))}
          </div>
        </div>

        <div className="mb-8">
          <p className="text-lg font-medium text-gray-700 dark:text-gray-200 transition-colors duration-500">
            {messages[messageIndex]}
            <span className="inline-block w-8 text-left text-blue-600 dark:text-blue-400 transition-colors duration-500">
              {'.'.repeat(dotCount)}
            </span>
          </p>
        </div>

        <div className="flex items-center justify-center space-x-2">
          <div
            className="w-2 h-2 bg-green-600 dark:bg-green-500 rounded-full animate-pulse transition-colors duration-500"
            style={{ animationDuration: '1s' }}
          ></div>
          <span className="text-xs font-medium text-gray-600 dark:text-gray-400 transition-colors duration-500">
            Blockchain network active
          </span>
        </div>
      </div>

      <style jsx>{`
        @keyframes pulse {
          0%,
          100% {
            opacity: 0.4;
            transform: scale(1);
          }
          50% {
            opacity: 1;
            transform: scale(1.1);
          }
        }
      `}</style>
    </div>
  );
};

export default RootLoader;
