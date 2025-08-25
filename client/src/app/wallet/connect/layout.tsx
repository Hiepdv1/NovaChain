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

  const { data, isLoading } = useWalletContext();

  useEffect(() => {
    if (data) {
      router.push('/');
    }
  }, [data, router]);

  if (isLoading) {
    return null;
  }

  return (
    <Fragment>
      {!isLoading && !data && (
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
