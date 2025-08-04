import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';
import { forwardRef, memo } from 'react';

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {}

const buttonVariants = cva('w-full text-white font-medium', {
  variants: {
    variant: {
      default: '',
      glass: 'glass-button magnetic-hover font-medium rounded-xl',
      secondary: 'secondary-button magnetic-hover font-semibold rounded-2xl',
      quantum: 'quantum-btn rounded-2xl',
    },
    size: {
      sm: 'py-2 px-3 text-sm',
      md: 'py-3 px-6 text-base',
      lg: 'py-4 px-8 text-lg',
    },

    defaultVariants: {
      variant: 'default',
      size: 'md',
    },
  },
});

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, children, ...props }, ref) => {
    return (
      <button
        ref={ref}
        {...props}
        className={cn(buttonVariants({ variant, size }), className)}
      >
        {children}
      </button>
    );
  },
);

Button.displayName = 'Button';

export default memo(Button);
