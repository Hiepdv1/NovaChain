import { Fragment } from 'react';

const Floating = () => {
  return (
    <Fragment>
      <div
        style={{
          animationDelay: '0s',
          background:
            'linear-gradient(135deg, rgba(255, 255, 255, 0.4), rgba(255, 255, 255, 0.1))',
        }}
        className="blur-[1px] animate-float absolute rounded-full w-32 h-32 top-10 left-10 opacity-30"
      />

      <div
        style={{
          animationDelay: '0s',
          background:
            'linear-gradient(135deg, rgba(255, 255, 255, 0.4), rgba(255, 255, 255, 0.1))',
        }}
        className="blur-[1px] animate-float absolute rounded-full w-20 h-20 top-1/3 right-20 opacity-30"
      />

      <div
        style={{
          animationDelay: '0s',
          background:
            'linear-gradient(135deg, rgba(255, 255, 255, 0.4), rgba(255, 255, 255, 0.1))',
        }}
        className="blur-[1px] animate-float absolute rounded-full w-24 h-24 bottom-20 left-1/4 opacity-30"
      />
    </Fragment>
  );
};

export default Floating;
