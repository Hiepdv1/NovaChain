import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export const IsValidNumber = (value: string): boolean => {
  if (!/^-?(?:\d+|\d+\.\d+)$/.test(value.trim())) return false;
  const num = Number(value);
  return !isNaN(num) && isFinite(num);
};

export const formatAddress = (addr: string) =>
  `${addr.slice(0, 6)}...${addr.slice(-4)}`;
