import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';
import { forwardRef, memo } from 'react';

export interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {}

const buttonVariants = cva('w-full text-white font-medium cursor-pointer', {
  variants: {
    variant: {
      default:
        'hover:bg-[#f1f5f9] hover:border-[#cbd5e1] hover:transform hover:-translate-y-0.5 font-medium text-gray-700 rounded-xl flex-1 bg-[#f8fafc] border-[2px] border-solid border-[#e2e8f0] transition-all duration-200 ease-in',
      elegant:
        'hover:bg-[linear-gradient(135deg,_#f43f5e_0%,_#ea580c_100%)] hover:shadow-[0_6px_20px_rgba(251,113,133,0.4)] hover:transform hover:-translate-y-0.5 bg-[linear-gradient(135deg,_#fb7185_0%,_#f97316_100%)] transition-all duration-200 ease-in shadow-[0_4px_15px_rgba(251,113,133,0.3)] font-medium rounded-xl flex flex-1 items-center justify-center',
      glass: 'glass-button magnetic-hover font-medium rounded-xl',
      secondary: 'secondary-button magnetic-hover font-semibold rounded-2xl',
      quantum: 'quantum-btn rounded-2xl',
    },
    size: {
      sm: 'py-2 px-4 text-sm',
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
        className={cn(buttonVariants({ variant, size }), className)}
        {...props}
      >
        {children}
      </button>
    );
  },
);

Button.displayName = 'Button';

export default memo(Button);
