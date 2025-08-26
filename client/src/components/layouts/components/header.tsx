import Button from '@/components/button';
import Input from '@/components/input';
import { ModeToggle } from '@/components/mode-toggle';

const Header = () => {
  return (
    <header className="shadow-xl dark:shadow-lg dark:ring-2 dark:ring-white/10 backdrop-blur-lg bg-gradient-glass dark:bg-gradient-glass-dark border-b border-white/20 dark:border-gray-700/50 sticky top-0 z-50">
      <div className="max-w-8xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <div className="flex items-center space-x-3">
            <div className="w-10 h-10 bg-gradient bg-gradient-primary rounded-xl flex items-center justify-center">
              <svg
                className="w-6 h-6 text-white"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>

            <div>
              <h1 className="text-sm font-bold bg-gradient-primary bg-clip-text text-transparent">
                CryptoChain
              </h1>
              <p className="text-[10px] text-gray-600 dark:text-gray-400">
                Explorer
              </p>
            </div>
          </div>

          {/* Search Bar */}
          <div className="hidden md:flex items-center flex-1 max-w-lg mx-8">
            <div className="relative w-full">
              <Input
                variant="levitating"
                inputSize="sm"
                className="dark:text-white pl-10 pr-4 dark:placeholder:text-white placeholder:text-black border-slate-400 font-normal py-3 rounded-xl !bg-gradient-glass text-xs text-black"
                id="search"
                name="s"
                placeholder="Search by address, tx hash, or block ...."
              />

              <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                <svg
                  className="w-4.5 h-4.5 text-gray-400"
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

              <button className="cursor-pointer absolute inset-y-0 right-0 pr-3 flex items-center">
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
            </div>
          </div>

          {/* Right Side */}
          <div className="flex items-center space-x-4">
            <button className="md:hidden p-2 bg-white rounded-full dark:bg-slate-600">
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
            </button>

            <ModeToggle className="text-xs cursor-pointer rounded-full w-9 h-9" />

            <Button
              variant="secondary"
              size="sm"
              className="text-sm rounded-lg"
            >
              Connect
            </Button>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
