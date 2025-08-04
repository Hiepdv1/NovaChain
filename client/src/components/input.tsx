import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';
import { forwardRef, memo } from 'react';
import React from 'react';

export interface InputProps
  extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'>,
    VariantProps<typeof inputVariants> {}

const inputVariants = cva('w-full text-white', {
  variants: {
    variant: {
      default: '',
      levitating:
        'rounded-2xl levitating-input text-slate-700 font-medium placeholder-transparent',
    },
    inputSize: {
      sm: 'px-4 pt-6 pb-2',
      md: 'px-6 pt-7 pb-3',
      lg: 'px-6 pt-8 pb-3',
    },
  },
  defaultVariants: {
    variant: 'default',
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
