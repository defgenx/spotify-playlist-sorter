import { NavLink } from 'react-router-dom';
import { Home, Music2, GitCompare } from 'lucide-react';
import { clsx } from 'clsx';

const navigation = [
  { name: 'Dashboard', href: '/', icon: Home },
  { name: 'Genres', href: '/genres', icon: Music2 },
  { name: 'Changes', href: '/changes', icon: GitCompare },
];

export function Sidebar() {
  return (
    <aside className="w-64 bg-spotify-black border-r border-spotify-gray min-h-[calc(100vh-73px)]">
      <nav className="p-4 space-y-2">
        {navigation.map((item) => (
          <NavLink
            key={item.name}
            to={item.href}
            className={({ isActive }) =>
              clsx(
                'flex items-center space-x-3 px-4 py-3 rounded-lg transition-colors',
                isActive
                  ? 'bg-spotify-gray text-white'
                  : 'text-spotify-lightgray hover:text-white hover:bg-spotify-darkgray'
              )
            }
          >
            <item.icon className="w-5 h-5" />
            <span className="font-medium">{item.name}</span>
          </NavLink>
        ))}
      </nav>
    </aside>
  );
}
