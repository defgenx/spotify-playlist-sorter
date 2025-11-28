import { Music, LogOut } from 'lucide-react';
import { useAuth } from '@/hooks/useAuth';
import { useUIStore } from '@/stores/uiStore';
import { Button } from '@/components/ui/Button';
import { Toggle } from '@/components/ui/Toggle';
import { Badge } from '@/components/ui/Badge';

export function Header() {
  const { user, logout } = useAuth();
  const { isDryRun, toggleDryRun } = useUIStore();

  return (
    <header className="bg-spotify-black border-b border-spotify-gray sticky top-0 z-50">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-3">
            <Music className="w-8 h-8 text-spotify-green" />
            <h1 className="text-2xl font-bold text-white">Spotify Playlist Sorter</h1>
          </div>

          <div className="flex items-center space-x-6">
            {user && (
              <>
                <div className="flex items-center space-x-4">
                  <Badge variant={isDryRun ? 'success' : 'danger'}>
                    {isDryRun ? 'Preview Mode' : 'Live Mode'}
                  </Badge>
                  <Toggle
                    checked={isDryRun}
                    onChange={toggleDryRun}
                    label="Dry Run"
                  />
                </div>

                <div className="flex items-center space-x-3">
                  {user.imageUrl && (
                    <img
                      src={user.imageUrl}
                      alt={user.displayName}
                      className="w-8 h-8 rounded-full"
                    />
                  )}
                  <span className="text-white text-sm">{user.displayName}</span>
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={logout}
                    className="gap-2"
                  >
                    <LogOut className="w-4 h-4" />
                    Logout
                  </Button>
                </div>
              </>
            )}
          </div>
        </div>
      </div>
    </header>
  );
}
