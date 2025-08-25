'use client';

import { ReactNode } from 'react';
import { Action, Toaster, toast as baseToast } from 'sonner';

const baseStyle =
  'flex justify-center items-center gap-3 px-5 py-4 rounded-xl text-base font-semibold shadow-lg transition-colors backdrop-blur-sm';

export const toast = {
  success: (message: string, icon?: ReactNode, action?: ReactNode | Action) =>
    baseToast(message, {
      action: action,
      icon: icon,
      className: `${baseStyle} border-none! text-white! bg-gradient-to-r! from-emerald-500! to-teal-600!`,
    }),

  error: (message: string, icon?: ReactNode, action?: ReactNode | Action) =>
    baseToast(message, {
      action: action,
      icon: icon,
      className: `${baseStyle} border-none! text-white! bg-gradient-to-r! from-red-500! to-pink-600!`,
    }),

  warning: (message: string, icon?: ReactNode, action?: ReactNode | Action) =>
    baseToast(message, {
      action: action,
      icon: icon,
      className: `${baseStyle} border-none! text-white! bg-gradient-to-r! from-amber-400! to-orange-500!`,
    }),

  info: (message: string, icon?: ReactNode, action?: ReactNode | Action) =>
    baseToast(message, {
      action: action,
      icon: icon,
      className: `${baseStyle} border-none! text-white! bg-gradient-to-r! from-sky-500! to-indigo-500!`,
    }),

  raw: baseToast,
};

export const GlobalToaster = () => (
  <Toaster
    position="top-right"
    toastOptions={{
      duration: 4000,
    }}
  />
);
