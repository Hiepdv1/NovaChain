'use client';

import { createContext, useContext } from 'react';
import { useWalletQuery } from '@/features/wallet/hook/useWalletQuery';

const WalletContext = createContext<ReturnType<typeof useWalletQuery> | null>(
  null,
);
export const WalletProvider = ({ children }: { children: React.ReactNode }) => {
  const walletQuery = useWalletQuery({
    retry: false,
    refetchOnWindowFocus: false,
    refetchOnReconnect: false,
    refetchOnMount: false,
  });

  return (
    <WalletContext.Provider value={walletQuery}>
      {children}
    </WalletContext.Provider>
  );
};

const useWalletContext = () => {
  const ctx = useContext(WalletContext);
  if (!ctx) {
    throw new Error('useWalletContext must be used inside WalletProvider');
  }
  return ctx;
};

export default useWalletContext;
