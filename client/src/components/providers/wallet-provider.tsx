'use client';

import { createContext, useContext } from 'react';
import { useWalletQuery } from '@/features/wallet/hook/useWalletQuery';
import RootLoader from '@/components/loading/root';

type walletContextProps = {
  wallet: ReturnType<typeof useWalletQuery> | null;
  refetch: (() => Promise<void>) | null;
};

const WalletContext = createContext<walletContextProps>({
  wallet: null,
  refetch: null,
});
export const WalletProvider = ({ children }: { children: React.ReactNode }) => {
  const walletQuery = useWalletQuery({
    retry: false,
    refetchOnWindowFocus: false,
    refetchOnReconnect: false,
    refetchOnMount: false,
  });

  const refetch = async () => {
    await Promise.all([walletQuery.refetch()]);
  };

  if (walletQuery.isLoading) {
    return <RootLoader />;
  }

  return (
    <WalletContext.Provider
      value={{
        wallet: walletQuery,
        refetch,
      }}
    >
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
