import { Fragment } from 'react';

const EnergyOrb = () => {
  return (
    <Fragment>
      <div
        style={{
          animationDelay: '0.5ms',
        }}
        className="absolute bg-spotlight animate-energy-pulse rounded-full w-20 h-20 opacity-30"
      />

      <div
        style={{
          animationDelay: '2.5ms',
        }}
        className="absolute bg-spotlight animate-energy-pulse rounded-full w-16 h-16 opacity-30"
      />

      <div
        style={{
          animationDelay: '4.5ms',
        }}
        className="absolute bg-spotlight animate-energy-pulse rounded-full w-24 h-24 opacity-30"
      />
    </Fragment>
  );
};

export default EnergyOrb;
