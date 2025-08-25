import { Fragment } from 'react';

const Particle = () => {
  return (
    <Fragment>
      <div
        className="opacity-0 absolute top-[5%] left-[20%] w-1.5 h-1.5 bg-particle-gradient rounded-full animate-particle-float"
        style={{ animationDelay: '0ms' }}
      ></div>
      <div
        className="opacity-0 absolute top-[120px] left-[83%] w-1.5 h-1.5 bg-particle-gradient rounded-full animate-particle-float"
        style={{ animationDelay: '500ms' }}
      ></div>
      <div
        className="opacity-0 absolute top-[28%] left-[7%] w-1.5 h-1.5 bg-particle-gradient rounded-full animate-particle-float"
        style={{ animationDelay: '900ms' }}
      ></div>
      <div
        className="opacity-0 absolute top-[220px] left-[65px] w-1.5 h-1.5 bg-particle-gradient rounded-full animate-particle-float"
        style={{ animationDelay: '200ms' }}
      ></div>
      <div
        className="opacity-0 absolute top-[62%] left-[52%] w-1.5 h-1.5 bg-particle-gradient rounded-full animate-particle-float"
        style={{ animationDelay: '650ms' }}
      ></div>

      <div
        className="opacity-0 absolute top-[48%] left-[45%] w-1.5 h-1.5 bg-particle-gradient rounded-full animate-particle-float"
        style={{ animationDelay: '300ms' }}
      ></div>
      <div
        className="opacity-0 absolute top-[55%] left-[50%] w-1.5 h-1.5 bg-particle-gradient rounded-full animate-particle-float"
        style={{ animationDelay: '600ms' }}
      ></div>
      <div
        className="opacity-0 absolute top-[50%] left-[58%] w-1.5 h-1.5 bg-particle-gradient rounded-full animate-particle-float"
        style={{ animationDelay: '150ms' }}
      ></div>
      <div
        className="opacity-0 absolute top-[46%] left-[52%] w-1.5 h-1.5 bg-particle-gradient rounded-full animate-particle-float"
        style={{ animationDelay: '950ms' }}
      ></div>

      <div
        className="opacity-0 absolute top-[42%] right-[6%] w-1.5 h-1.5 bg-particle-gradient rounded-full animate-particle-float"
        style={{ animationDelay: '100ms' }}
      ></div>
      <div
        className="opacity-0 absolute top-[165px] right-[140px] w-1.5 h-1.5 bg-particle-gradient rounded-full animate-particle-float"
        style={{ animationDelay: '750ms' }}
      ></div>
      <div
        className="opacity-0 absolute top-[78%] right-[23%] w-1.5 h-1.5 bg-particle-gradient rounded-full animate-particle-float"
        style={{ animationDelay: '350ms' }}
      ></div>
      <div
        className="opacity-0 absolute top-[12%] right-[75px] w-1.5 h-1.5 bg-particle-gradient rounded-full animate-particle-float"
        style={{ animationDelay: '1150ms' }}
      ></div>
      <div
        className="opacity-0 absolute top-[88%] right-[12%] w-1.5 h-1.5 bg-particle-gradient rounded-full animate-particle-float"
        style={{ animationDelay: '550ms' }}
      ></div>
    </Fragment>
  );
};

export default Particle;
