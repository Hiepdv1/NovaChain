'use client';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

interface NavLink {
  id: number;
  title: string;
  to: string;
  active: string[];
}

const NavList: NavLink[] = [
  {
    id: 1,
    title: 'Dashboard',
    to: '/',
    active: ['/'],
  },
  {
    id: 2,
    title: 'Blocks',
    to: '/blocks',
    active: ['/block'],
  },
  {
    id: 3,
    title: 'Transactions',
    to: '/txs',
    active: ['/txs', '/txs/send'],
  },
  {
    id: 4,
    title: 'Pending',
    to: '/txs/pending',
    active: ['/txs/pending'],
  },
  {
    id: 5,
    title: 'Miners',
    to: '/miners',
    active: ['/miners'],
  },
];

const NavBar = () => {
  const pathname = usePathname();

  return (
    <div className="overflow-x-auto scrollbar-hide mb-8 select-none glass-card dark:border-white/5 dark:bg-transparent dark:!bg-[linear-gradient(135deg,_rgba(255,255,255,0.1)_0%,_rgba(255,255,255,0.05)_100%)] p-1 rounded-xl">
      <nav aria-label="Main navigation" className="flex gap-x-1">
        {NavList.map((navItem) => {
          const isActive = navItem.active.includes(pathname);

          const classActive = isActive
            ? 'bg-white text-gray-950 dark:bg-gray-700 dark:text-white'
            : 'dark:hover:text-white';

          return (
            <Link
              className={`${classActive} hover:text-gray-950 text-gray-600 dark:text-gray-400 transition-all duration-200 rounded-lg font-medium text-xs px-4 py-3`}
              key={navItem.id}
              href={navItem.to}
            >
              <span>{navItem.title}</span>
            </Link>
          );
        })}
      </nav>
    </div>
  );
};

export default NavBar;
