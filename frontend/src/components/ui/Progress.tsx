import { HTMLAttributes, forwardRef } from 'react';
import { clsx } from 'clsx';

interface ProgressProps extends HTMLAttributes<HTMLDivElement> {
  value: number;
  max?: number;
  showLabel?: boolean;
}

export const Progress = forwardRef<HTMLDivElement, ProgressProps>(
  ({ className, value, max = 100, showLabel = false, ...props }, ref) => {
    const percentage = Math.min(Math.max((value / max) * 100, 0), 100);

    return (
      <div ref={ref} className={clsx('w-full', className)} {...props}>
        <div className="flex items-center justify-between mb-1">
          {showLabel && (
            <span className="text-sm text-spotify-lightgray">
              {Math.round(percentage)}%
            </span>
          )}
        </div>
        <div className="w-full bg-spotify-gray rounded-full h-2 overflow-hidden">
          <div
            className="h-full bg-spotify-green transition-all duration-300 ease-in-out"
            style={{ width: `${percentage}%` }}
          />
        </div>
      </div>
    );
  }
);

Progress.displayName = 'Progress';
