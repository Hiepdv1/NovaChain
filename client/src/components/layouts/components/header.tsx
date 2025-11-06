'use client';

import Button from '@/components/button';
import { toast } from '@/components/globalToaster';
import Input from '@/components/input';
import { ModeToggle } from '@/components/mode-toggle';
import useWalletContext from '@/components/providers/wallet-provider';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { useDisconnectWalletMutation } from '@/features/wallet/hook/useWalletMutation';
import { DelWalletPool, GetWalletPool } from '@/lib/db/wallet.index';
import { formatAddress } from '@/lib/utils';
import { StoredWallet } from '@/shared/types/wallet';
import {
  Activity,
  Check,
  ChevronDown,
  Copy,
  ExternalLink,
  LogOut,
  Send,
  Settings,
  Wallet,
  X,
} from 'lucide-react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { FormEvent, useEffect, useState } from 'react';

interface MenuItem {
  id: number;
  title: string;
  icon: React.ElementType;
  path: string;
}

const MenuList: MenuItem[] = [
  {
    id: 1,
    title: 'View on Explorer',
    icon: ExternalLink,
    path: '/',
  },
  {
    id: 2,
    title: 'Transaction History',
    icon: Activity,
    path: '/txs',
  },
  {
    id: 3,
    title: 'Send Transaction',
    icon: Send,
    path: '/txs/send',
  },
  {
    id: 4,
    title: 'Wallet Settings',
    icon: Settings,
    path: '/settings',
  },
];

const Header = () => {
  const { wallet: walletQuery, refetch } = useWalletContext();
  const [wallet, setWallet] = useState<StoredWallet | null>(null);
  const [copied, setCopied] = useState(false);
  const { mutate: disconnectWallet } = useDisconnectWalletMutation();
  const searchParams = useSearchParams();
  const query = searchParams.get('result_search');
  const router = useRouter();
  const [activeSearch, setActiveSearch] = useState(false);

  useEffect(() => {
    GetWalletPool().then((ws) => {
      if (ws.length > 0) {
        setWallet(ws[0]);
      }
    });
  }, []);

  const onCopy = async () => {
    if (!walletQuery || !walletQuery.data || !wallet) return;

    setCopied(true);

    await navigator.clipboard.writeText(wallet.address);

    toast.success('🎉 Address copied!');
    setTimeout(() => {
      setCopied(false);
    }, 3000);
  };

  const onDisconnect = () => {
    disconnectWallet(null, {
      onSuccess: async () => {
        await DelWalletPool();
        await refetch?.();
        setWallet(null);
      },
    });
  };

  const onConnect = () => {
    router.push('/wallet/connect');
    router.refresh();
  };

  const onNavLink = async (item: MenuItem) => {
    await refetch?.();
    router.push(item.path);
    router.refresh();
  };

  const onSubmit = (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const formData = new FormData(e.currentTarget);
    const data = Object.fromEntries(formData) as { result_search: string };

    router.push(`/search?result_search=${data.result_search}`);
  };

  return (
    <header className="shadow-xl dark:shadow-lg dark:ring-2 dark:ring-white/10 backdrop-blur-lg bg-gradient-glass dark:bg-gradient-glass-dark border-b border-white/20 dark:border-gray-700/50 sticky top-0 z-50">
      <div className="max-w-8xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <div className="flex items-center space-x-3">
            <Link href="/">
              <div className="w-12 h-12 bg-gradient bg-gradient-primary rounded-xl flex items-center justify-center">
                <svg
                  className="relative z-10 w-10 h-10 text-white animate-pulse"
                  fill="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-1 17.93c-3.94-.49-7-3.85-7-7.93 0-.62.08-1.21.21-1.79L9 15v1c0 1.1.9 2 2 2v1.93zm6.9-2.54c-.26-.81-1-1.39-1.9-1.39h-1v-3c0-.55-.45-1-1-1H8v-2h2c.55 0 1-.45 1-1V7h2c1.1 0 2-.9 2-2v-.41c2.93 1.19 5 4.06 5 7.41 0 2.08-.8 3.97-2.1 5.39z" />
                </svg>
              </div>
            </Link>

            <div>
              <h1 className="text-xl font-bold bg-gradient-primary bg-clip-text text-transparent">
                CryptoChain
              </h1>
              <p className="text-sm text-gray-600 dark:text-gray-400">
                Explorer
              </p>
            </div>
          </div>

          {/* Search Bar */}
          <div className="hidden md:flex items-center flex-1 max-w-lg mx-8">
            <form onSubmit={onSubmit} className="relative w-full">
              <Input
                variant="levitating"
                inputSize="sm"
                className="dark:text-white pl-10 pr-8 dark:!bg-transparent dark:placeholder:text-white placeholder:text-black border-slate-400 font-normal py-3 rounded-xl !bg-gradient-glass text-sm text-black"
                id="search"
                name="result_search"
                defaultValue={query || ''}
                placeholder="Search by address, tx hash, or block ...."
              />

              <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <svg
                  className="w-5 h-5 text-gray-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                  ></path>
                </svg>
              </div>

              <button
                type="submit"
                className="cursor-pointer absolute inset-y-0 right-0 pr-3 flex items-center"
              >
                <svg
                  className="w-5 h-5 text-primary-color dark:text-primary-400 hover:text-primary-700 dark:hover:text-primary-300"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M13 7l5 5-5 5M6 12h12"
                  ></path>
                </svg>
              </button>
            </form>
          </div>

          {/* Right Side */}
          <div className="relative flex items-center max-sm:space-x-3 space-x-4">
            <button
              onClick={() => setActiveSearch((prev) => !prev)}
              className="md:hidden p-2 bg-white rounded-full dark:bg-slate-600"
            >
              {!activeSearch && (
                <svg
                  className="w-5 h-5 text-gray-700 dark:text-gray-300"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                  ></path>
                </svg>
              )}

              {activeSearch && (
                <X className="w-5 h-5 text-gray-700 dark:text-gray-300" />
              )}
            </button>

            {activeSearch && (
              <form
                onSubmit={onSubmit}
                className="animate-cascase-fade absolute md:hidden -bottom-[calc(100%+8px)] -right-7 w-screen"
              >
                <Input
                  variant="levitating"
                  inputSize="sm"
                  className="dark:text-white z-50 px-4 dark:placeholder:text-white placeholder:text-black border-slate-400 font-normal py-3 rounded-xl !bg-gradient-glass text-sm text-black"
                  id="search"
                  name="result_search"
                  type="search"
                  defaultValue={query || ''}
                  placeholder="Search by address, tx hash, or block ...."
                />
              </form>
            )}

            <ModeToggle className="text-xs cursor-pointer rounded-full w-9 h-9" />

            {wallet ? (
              <div className="outline-none select-none">
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <button className="outline-none cursor-pointer flex items-center max-sm:space-x-1 space-x-3 px-4 py-2 bg-gradient-to-r from-blue-50 to-purple-50 dark:from-blue-900/20 dark:to-purple-900/20 hover:from-blue-100 hover:to-purple-100 dark:hover:from-blue-900/30 dark:hover:to-purple-900/30 border border-blue-200/50 dark:border-blue-700/50 rounded-xl transition-all duration-200 group">
                      <div className="relative !mr-0">
                        <div className="w-8 h-8 bg-gradient-to-br from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
                          <Wallet className="w-4 h-4 text-white" />
                        </div>
                        <div className="absolute -top-1 -right-1 w-3 h-3 bg-green-500 border-2 border-white dark:border-gray-900 rounded-full"></div>
                      </div>

                      <div className="text-left hidden md:block">
                        <div className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                          {formatAddress(wallet.address)}
                        </div>
                        <div className="flex items-center space-x-1 text-xs text-gray-500 dark:text-gray-400">
                          <span>
                            {(walletQuery?.data?.data.Balance &&
                              parseFloat(
                                walletQuery.data.data.Balance,
                              ).toString()) ||
                              'N/A'}{' '}
                            CCC
                          </span>
                          <span>•</span>
                        </div>
                      </div>

                      <ChevronDown className="ml-3 max-sm:hidden w-4 h-4 text-gray-400 group-hover:text-gray-600 dark:group-hover:text-gray-300 transition-colors duration-200" />
                    </button>
                  </DropdownMenuTrigger>

                  <DropdownMenuContent className="w-80 bg-white dark:bg-gray-900 border border-gray-200 dark:border-gray-800 shadow-xl">
                    <DropdownMenuLabel className="px-4 py-3">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-3">
                          <div className="w-10 h-10 bg-gradient-to-br from-blue-600 to-purple-600 rounded-lg flex items-center justify-center">
                            <Wallet className="w-5 h-5 text-white" />
                          </div>
                          <div>
                            <div className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                              Connected Wallet
                            </div>
                            <div className="flex items-center space-x-1 text-xs text-gray-500 dark:text-gray-400">
                              <div className="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
                              <span>Active • CCC</span>
                            </div>
                          </div>
                        </div>
                      </div>
                    </DropdownMenuLabel>

                    <DropdownMenuSeparator className="bg-gray-200 dark:bg-gray-800" />

                    <div className="px-4 py-3 bg-gray-50 dark:bg-gray-800/50">
                      <div className="grid grid-cols-2 gap-4">
                        <div>
                          <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">
                            Balance
                          </div>
                          <div className="text-sm font-bold text-gray-900 dark:text-gray-100">
                            {(walletQuery?.data?.data.Balance &&
                              parseFloat(
                                walletQuery.data.data.Balance,
                              ).toString()) ||
                              'N/A'}
                          </div>
                        </div>
                      </div>
                    </div>

                    <div className="px-4 py-2">
                      <div className="text-xs text-gray-500 dark:text-gray-400 mb-1">
                        Wallet Address
                      </div>
                      <div className="flex items-center justify-between p-2 bg-gray-100 dark:bg-gray-800 rounded-lg">
                        <code className="text-sm font-mono text-gray-900 dark:text-gray-100">
                          {formatAddress(wallet.address)}
                        </code>
                        <button
                          disabled={copied}
                          onClick={onCopy}
                          className="cursor-pointer p-1 hover:bg-gray-200 dark:hover:bg-gray-700 rounded transition-colors duration-200"
                        >
                          {copied ? (
                            <Check className="w-3 h-3 text-gray-500 dark:text-gray-400" />
                          ) : (
                            <Copy className="w-3 h-3 text-gray-500 dark:text-gray-400" />
                          )}
                        </button>
                      </div>
                    </div>

                    <DropdownMenuSeparator className="bg-gray-200 dark:bg-gray-800" />

                    {MenuList.map((navItem) => {
                      const Icon = navItem.icon;

                      return (
                        <DropdownMenuItem
                          onClick={() => onNavLink(navItem)}
                          key={navItem.id}
                          className="px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-800 cursor-pointer"
                        >
                          <Icon className="w-4 h-4 mr-3 text-gray-500 dark:text-gray-400" />
                          <span className="text-xs text-gray-700 dark:text-gray-300">
                            {navItem.title}
                          </span>
                        </DropdownMenuItem>
                      );
                    })}

                    <DropdownMenuSeparator className="bg-gray-200 dark:bg-gray-800" />

                    <DropdownMenuItem
                      onClick={onDisconnect}
                      className="px-4 py-3 hover:bg-red-50 dark:hover:bg-red-900/20 cursor-pointer text-red-600 dark:text-red-400"
                    >
                      <LogOut className="w-4 h-4 mr-3" />
                      <span className="text-xs">Disconnect Wallet</span>
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            ) : (
              <Button
                variant="secondary"
                size="sm"
                className="rounded-lg flex items-center"
                onClick={onConnect}
              >
                <Wallet className="w-4 h-4 mr-2" />
                <span className="text-sm">Connect</span>
              </Button>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
