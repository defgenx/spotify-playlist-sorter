import { Plus, Minus, FolderPlus } from 'lucide-react';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';

interface ChangeSummaryProps {
  tracksToAdd: number;
  tracksToRemove: number;
  playlistsToCreate: number;
}

export function ChangeSummary({ tracksToAdd, tracksToRemove, playlistsToCreate }: ChangeSummaryProps) {
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      <Card variant="highlighted">
        <CardHeader>
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-green-600 rounded-lg">
              <Plus className="w-5 h-5 text-white" />
            </div>
            <CardTitle>Tracks to Add</CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          <p className="text-3xl font-bold text-white">{tracksToAdd}</p>
        </CardContent>
      </Card>

      <Card variant="highlighted">
        <CardHeader>
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-red-600 rounded-lg">
              <Minus className="w-5 h-5 text-white" />
            </div>
            <CardTitle>Tracks to Remove</CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          <p className="text-3xl font-bold text-white">{tracksToRemove}</p>
        </CardContent>
      </Card>

      <Card variant="highlighted">
        <CardHeader>
          <div className="flex items-center space-x-3">
            <div className="p-2 bg-blue-600 rounded-lg">
              <FolderPlus className="w-5 h-5 text-white" />
            </div>
            <CardTitle>Playlists to Create</CardTitle>
          </div>
        </CardHeader>
        <CardContent>
          <p className="text-3xl font-bold text-white">{playlistsToCreate}</p>
        </CardContent>
      </Card>
    </div>
  );
}
