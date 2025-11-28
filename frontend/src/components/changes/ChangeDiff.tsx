import { useState } from 'react';
import type { TrackMove as TrackMoveType } from '@/lib/types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card';
import { TrackMove } from './TrackMove';
import { Button } from '@/components/ui/Button';

interface ChangeDiffProps {
  tracksToAdd: TrackMoveType[];
  tracksToRemove: TrackMoveType[];
}

export function ChangeDiff({ tracksToAdd, tracksToRemove }: ChangeDiffProps) {
  const [showAdditions, setShowAdditions] = useState(true);
  const [showRemovals, setShowRemovals] = useState(true);

  return (
    <div className="space-y-6">
      {tracksToAdd.length > 0 && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="text-green-500">
                Tracks to Add ({tracksToAdd.length})
              </CardTitle>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setShowAdditions(!showAdditions)}
              >
                {showAdditions ? 'Hide' : 'Show'}
              </Button>
            </div>
          </CardHeader>
          {showAdditions && (
            <CardContent className="space-y-3">
              {tracksToAdd.map((move, idx) => (
                <TrackMove key={idx} move={move} variant="add" />
              ))}
            </CardContent>
          )}
        </Card>
      )}

      {tracksToRemove.length > 0 && (
        <Card>
          <CardHeader>
            <div className="flex items-center justify-between">
              <CardTitle className="text-red-500">
                Tracks to Remove ({tracksToRemove.length})
              </CardTitle>
              <Button
                variant="ghost"
                size="sm"
                onClick={() => setShowRemovals(!showRemovals)}
              >
                {showRemovals ? 'Hide' : 'Show'}
              </Button>
            </div>
          </CardHeader>
          {showRemovals && (
            <CardContent className="space-y-3">
              {tracksToRemove.map((move, idx) => (
                <TrackMove key={idx} move={move} variant="remove" />
              ))}
            </CardContent>
          )}
        </Card>
      )}
    </div>
  );
}
