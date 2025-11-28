import { Progress } from '@/components/ui/Progress';

interface ProgressBarProps {
  current: number;
  total: number;
  message?: string;
}

export function ProgressBar({ current, total, message }: ProgressBarProps) {
  const percentage = total > 0 ? (current / total) * 100 : 0;

  return (
    <div className="space-y-2">
      {message && (
        <p className="text-sm text-spotify-lightgray">{message}</p>
      )}
      <Progress value={percentage} max={100} showLabel />
      <p className="text-xs text-spotify-lightgray text-right">
        {current} / {total}
      </p>
    </div>
  );
}
