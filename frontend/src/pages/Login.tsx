import { Music } from 'lucide-react';
import { useAuth } from '@/hooks/useAuth';
import { Button } from '@/components/ui/Button';

export function Login() {
  const { login } = useAuth();

  return (
    <div className="min-h-screen bg-spotify-black flex items-center justify-center px-4">
      <div className="max-w-md w-full text-center space-y-8">
        <div className="space-y-4">
          <div className="flex justify-center">
            <div className="p-6 bg-spotify-green rounded-full">
              <Music className="w-16 h-16 text-white" />
            </div>
          </div>
          <h1 className="text-4xl font-bold text-white">
            Spotify Playlist Sorter
          </h1>
          <p className="text-spotify-lightgray text-lg">
            Automatically organize your Spotify library by genre
          </p>
        </div>

        <div className="space-y-4">
          <Button
            onClick={login}
            size="lg"
            className="w-full"
          >
            Login with Spotify
          </Button>
          <p className="text-xs text-spotify-lightgray">
            By logging in, you agree to grant access to your Spotify library
          </p>
        </div>

        <div className="pt-8 space-y-3 text-sm text-spotify-lightgray">
          <div className="flex items-center justify-center space-x-2">
            <CheckIcon />
            <span>Analyze your music library</span>
          </div>
          <div className="flex items-center justify-center space-x-2">
            <CheckIcon />
            <span>Sort tracks by genre automatically</span>
          </div>
          <div className="flex items-center justify-center space-x-2">
            <CheckIcon />
            <span>Preview changes before applying</span>
          </div>
        </div>
      </div>
    </div>
  );
}

function CheckIcon() {
  return (
    <svg
      className="w-5 h-5 text-spotify-green"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
    >
      <path
        strokeLinecap="round"
        strokeLinejoin="round"
        strokeWidth={2}
        d="M5 13l4 4L19 7"
      />
    </svg>
  );
}
