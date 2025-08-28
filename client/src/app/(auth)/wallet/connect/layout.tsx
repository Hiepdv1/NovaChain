'use client';

import useWalletContext from '@/components/providers/wallet-provider';
import EnergyOrb from '@/features/wallet/components/energy-orb';
import Floating from '@/features/wallet/components/floating';
import Particle from '@/features/wallet/components/particle';
import { useRouter } from 'next/navigation';
import { Fragment, useEffect } from 'react';

const WalletLayout = ({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) => {
  const router = useRouter();

  const { isLoading, error } = useWalletContext();

  useEffect(() => {
    if (error && error?.statusCode !== 401) {
      router.push('/');
    }
  }, [router, error]);

  return (
    <Fragment>
      {!isLoading && (
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
