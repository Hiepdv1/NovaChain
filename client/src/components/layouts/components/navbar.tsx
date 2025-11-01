'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { motion } from 'framer-motion';

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
    active: ['/blocks', '/blocks/[hash]'],
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

  const isMatch = (pathname: string, activeList: string[]) => {
    return activeList.some((pattern) => {
      if (pattern.includes('[')) {
        const base = pattern.split('/[')[0];
        return pathname.startsWith(base + '/');
      }
      return pathname === pattern;
    });
  };

  return (
    <div className="mb-8 select-none">
      <nav
        aria-label="Main navigation"
        className="overflow-x-auto scrollbar-hide relative flex gap-x-2 bg-gradient-to-r from-white/60 to-white/40 dark:from-gray-800/50 dark:to-gray-700/40 backdrop-blur-xl border border-gray-200/50 dark:border-gray-700/50 rounded-2xl p-2 shadow-sm"
      >
        {NavList.map((navItem) => {
          const isActive = isMatch(pathname, navItem.active);

          return (
            <Link
              key={navItem.id}
              href={navItem.to}
              className={`relative px-5 py-2.5 rounded-xl font-medium text-sm transition-all duration-200 whitespace-nowrap ${
                isActive
                  ? 'text-white dark:text-white'
                  : 'text-gray-600 dark:text-gray-400 hover:text-gray-900 dark:hover:text-gray-200'
              }`}
            >
              {isActive && (
                <motion.div
                  layoutId="active-pill"
                  className="absolute inset-0 rounded-xl bg-gradient-primary shadow-md"
                  transition={{ type: 'spring', stiffness: 350, damping: 25 }}
                />
              )}
              <span className="relative z-10">{navItem.title}</span>
            </Link>
          );
        })}
      </nav>
    </div>
  );
};

export default NavBar;
