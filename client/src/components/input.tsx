import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/shared/utils/class';
import { forwardRef, memo } from 'react';
import React from 'react';

export interface InputProps
  extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'>,
    VariantProps<typeof inputVariants> {}

const inputVariants = cva('w-full text-white', {
  variants: {
    variant: {
      levitating:
        'rounded-2xl levitating-input text-slate-700 font-medium placeholder-transparent',
      elegant:
        'focus:outline-none not-placeholder-shown:border-[#f97316] focus:border-[#f97316] focus:shadow-[0,_0,_0,_3px,_rgba(251,113,133,0.1)] text-gray-800 text-lg py-4 px-4  rounded-xl w-full border-gray-300 border-solid border-[2px] transition-all duration-300 ease-in-out w-full px-4 py-4 rounded-xl text-gray-800 placeholder-gray-400 text-lg',
    },
    inputSize: {
      sm: 'px-4 pt-6 pb-2',
      md: 'px-6 pt-7 pb-3',
      lg: 'px-6 pt-8 pb-3',
    },
  },
  defaultVariants: {
    variant: 'levitating',
    inputSize: 'md',
  },
});

const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ className, variant, inputSize, ...props }, ref) => {
    return (
      <input
        ref={ref}
        {...props}
        className={cn(inputVariants({ variant, inputSize }), className)}
      />
    );
  },
);

Input.displayName = 'Input';

export default memo(Input);
