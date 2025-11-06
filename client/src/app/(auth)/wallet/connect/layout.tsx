'use client';

import useWalletContext from '@/components/providers/wallet-provider';
import EnergyOrb from '@/features/wallet/components/connect/energy-orb';
import Floating from '@/features/wallet/components/connect/floating';
import Particle from '@/features/wallet/components/connect/particle';
import { useRouter } from 'next/navigation';
import { Fragment, useEffect } from 'react';

const WalletLayout = ({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) => {
  const router = useRouter();

  const { wallet } = useWalletContext();

  useEffect(() => {
    if (wallet?.error && wallet.error?.statusCode !== 401) {
      router.push('/');
    }
  }, [router, wallet?.error]);

  return (
    <Fragment>
      {!wallet?.isLoading && (
        <div className="min-h-screen flex items-center justify-center relative overflow-hidden bg-wallet-gradient bg-[400%,400%]! animate-gradient-shift ">
          <Floating />
          <Particle />
          <EnergyOrb />
          {children}
        </div>
      )}
    </Fragment>
  );
};

export default WalletLayout;
