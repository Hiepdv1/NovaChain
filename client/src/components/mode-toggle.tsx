'use client';

import * as React from 'react';
import { Moon, Sun } from 'lucide-react';
import { useTheme } from 'next-themes';

import { Button } from '@/components/ui/button';

export function ModeToggle(props: { className?: string }) {
  const { setTheme } = useTheme();

  const onModeToggle = () => {
    const mode = localStorage.getItem('theme');

    if (mode === 'dark') {
      setTheme('light');
    } else if (mode === 'light') {
      setTheme('dark');
    } else {
      if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
        setTheme('light');
      } else {
        setTheme('dark');
      }
    }
  };

  return (
    <Button
      onClick={onModeToggle}
      className={props.className}
      variant="outline"
      size="icon"
    >
      <Sun className="h-[1.2rem] w-[1.2rem] scale-100 rotate-0 transition-all dark:scale-0 dark:-rotate-90" />
      <Moon className="absolute h-[1.2rem] w-[1.2rem] scale-0 rotate-90 transition-all dark:scale-100 dark:rotate-0" />
    </Button>
  );
}
