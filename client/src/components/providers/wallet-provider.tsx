'use client';

import { createContext, useContext } from 'react';
import { useWalletQuery } from '@/features/wallet/hook/useWalletQuery';
import RootLoader from '@/components/loading/root';

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

  if (walletQuery.isLoading) {
    return <RootLoader />;
  }

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
