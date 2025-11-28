import { HTMLAttributes, forwardRef } from 'react';
import { clsx } from 'clsx';

interface BadgeProps extends HTMLAttributes<HTMLSpanElement> {
  variant?: 'default' | 'success' | 'warning' | 'danger' | 'info';
}

export const Badge = forwardRef<HTMLSpanElement, BadgeProps>(
  ({ className, variant = 'default', children, ...props }, ref) => {
    const variants = {
      default: 'bg-spotify-gray text-white',
      success: 'bg-spotify-green text-white',
      warning: 'bg-yellow-500 text-black',
      danger: 'bg-red-600 text-white',
      info: 'bg-blue-600 text-white',
    };

    return (
      <span
        ref={ref}
        className={clsx(
          'inline-flex items-center px-3 py-1 rounded-full text-xs font-medium',
          variants[variant],
          className
        )}
        {...props}
      >
        {children}
      </span>
    );
  }
);

Badge.displayName = 'Badge';
