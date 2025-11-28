import { useEffect } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Loader2 } from 'lucide-react';
import { api } from '@/lib/api';
import { useAuth } from '@/hooks/useAuth';

export function Callback() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const { checkAuth } = useAuth();

  useEffect(() => {
    const token = searchParams.get('token');
    const error = searchParams.get('error');

    if (error) {
      console.error('Auth error:', error);
      navigate('/login');
      return;
    }

    if (!token) {
      navigate('/login');
      return;
    }

    handleCallback(token);
  }, [searchParams, navigate]);

  const handleCallback = async (token: string) => {
    try {
      await api.completeLogin(token);
      await checkAuth();
      navigate('/');
    } catch (error) {
      console.error('Callback failed:', error);
      navigate('/login');
    }
  };

  return (
    <div className="min-h-screen bg-spotify-black flex items-center justify-center">
      <div className="text-center space-y-4">
        <Loader2 className="w-12 h-12 text-spotify-green animate-spin mx-auto" />
        <p className="text-white text-lg">Logging you in...</p>
      </div>
    </div>
  );
}
