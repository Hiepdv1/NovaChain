'use client';

import { createContext, useContext, useState } from 'react';
import { useRouter } from 'next/navigation';

type HistoryEntry = {
  key: string;
  title?: string;
  path: string;
};

type HistoryContextType = {
  history: HistoryEntry[];
  current: HistoryEntry | null;
  push: (entry: Omit<HistoryEntry, 'key'>) => void;
  back: () => void;
  clear: () => void;
};

const HistoryContext = createContext<HistoryContextType | undefined>(undefined);

export const HistoryProvider = ({
  children,
}: {
  children: React.ReactNode;
}) => {
  const router = useRouter();
  const [history, setHistory] = useState<HistoryEntry[]>([]);

  const push = (entry: Omit<HistoryEntry, 'key'>) => {
    setHistory((prev) => [...prev, { ...entry, key: crypto.randomUUID() }]);
    router.push(entry.path);
  };

  const back = () => {
    setHistory((prev) => {
      if (prev.length <= 1) return prev;
      const newHistory = [...prev];
      newHistory.pop();
      const last = newHistory[newHistory.length - 1];
      if (last) router.push(last.path);
      return newHistory;
    });
  };

  const clear = () => setHistory([]);

  const current = history[history.length - 1] || null;

  return (
    <HistoryContext.Provider value={{ history, current, push, back, clear }}>
      {children}
    </HistoryContext.Provider>
  );
};

export const useHistoryContext = () => {
  const context = useContext(HistoryContext);
  if (!context)
    throw new Error('useHistoryContext must be used within a HistoryProvider');
  return context;
};
