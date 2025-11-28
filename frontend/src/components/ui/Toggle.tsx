import { InputHTMLAttributes, forwardRef } from 'react';
import { clsx } from 'clsx';

interface ToggleProps extends Omit<InputHTMLAttributes<HTMLInputElement>, 'type'> {
  label?: string;
}

export const Toggle = forwardRef<HTMLInputElement, ToggleProps>(
  ({ className, label, checked, ...props }, ref) => {
    return (
      <label className={clsx('inline-flex items-center cursor-pointer', className)}>
        <div className="relative">
          <input
            ref={ref}
            type="checkbox"
            className="sr-only peer"
            checked={checked}
            {...props}
          />
          <div className="w-11 h-6 bg-spotify-gray rounded-full peer peer-checked:bg-spotify-green transition-colors" />
          <div className="absolute left-1 top-1 w-4 h-4 bg-white rounded-full transition-transform peer-checked:translate-x-5" />
        </div>
        {label && (
          <span className="ml-3 text-sm font-medium text-white">{label}</span>
        )}
      </label>
    );
  }
);

Toggle.displayName = 'Toggle';
