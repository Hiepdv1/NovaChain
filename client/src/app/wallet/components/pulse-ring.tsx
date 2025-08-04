import { Fragment } from 'react';

const PulseRing = () => {
  return (
    <Fragment>
      <div className="opacity-0 w-24 h-24 top-0 left-0 absolute border-2 border-solid border-[rgba(255,255,255,0.3)] rounded-full animate-pulse-ring" />

      <div
        style={{
          animationDelay: '1s',
        }}
        className="opacity-0 w-24 h-24 top-0 left-0 absolute border-2 border-solid border-[rgba(255,255,255,0.3)] rounded-full animate-pulse-ring"
      />

      <div className="opacity-0 w-24 h-24 top-0 left-0 absolute border-2 border-solid border-[rgba(255,255,255,0.3)] rounded-full animate-pulse-ring" />

      <div
        style={{
          animationDelay: '2s',
        }}
        className="opacity-0 w-24 h-24 top-0 left-0 absolute border-2 border-solid border-[rgba(255,255,255,0.3)] rounded-full animate-pulse-ring"
      />

      <div
        style={{
          animationDelay: '1s',
        }}
        className="opacity-0 w-24 h-24 top-0 left-0 absolute border-2 border-solid border-[rgba(255,255,255,0.3)] rounded-full animate-pulse-ring"
      />

      <div
        style={{
          animationDelay: '2s',
        }}
        className="opacity-0 w-24 h-24 top-0 left-0 absolute border-2 border-solid border-[rgba(255,255,255,0.3)] rounded-full animate-pulse-ring"
      />
    </Fragment>
  );
};

export default PulseRing;
