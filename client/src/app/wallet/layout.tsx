'use client';

import EnergyOrb from './components/energy-orb';
import Floating from './components/floating';
import Particle from './components/particle';

const WalletLayout = ({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) => {
  return (
    <div className="min-h-screen flex items-center justify-center p-4 relative overflow-hidden bg-wallet-gradient bg-[400%,400%]! animate-gradient-shift ">
      <Floating />
      <Particle />
      <EnergyOrb />
      {children}
    </div>
  );
};

export default WalletLayout;
