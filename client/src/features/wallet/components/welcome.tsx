import { memo } from 'react';
import PulseRing from './pulse-ring';
import { ModalName } from '../pages/page';
import Button from '@/components/button';

interface WelcomeProps {
  onSwitchModal(modalName: ModalName): void;
}

const Welcome = ({ onSwitchModal }: WelcomeProps) => {
  const onShowCreateWalletModal = () => {
    onSwitchModal('CreateWallet');
  };

  const onShowImportWalletModal = () => {
    onSwitchModal('ImportWallet');
  };

  return (
    <div className="glass-card rounded-3xl p-8">
      <div className="mb-8">
        <div className="relative w-24 h-24 mx-auto mb-6">
          <div className="w-24 h-24 bg-gradient-to-br from-blue-500 via-cyan-300 to-green-500 rounded-2xl flex items-center justify-center animate-bounce-in">
            <svg
              className="w-12 h-12 text-white"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path d="M12,1L3,5V11C3,16.55 6.84,21.74 12,23C17.16,21.74 21,16.55 21,11V5L12,1M12,7C13.66,7 15,8.34 15,10V11H16V17H8V11H9V10C9,8.34 10.34,7 12,7M12,8.5C11.17,8.5 10.5,9.17 10.5,10V11H13.5V10C13.5,9.17 12.83,8.5 12,8.5Z"></path>
            </svg>
          </div>
          <PulseRing />
        </div>

        <div className="text-center">
          <h1 className=" text-4xl font-bold text-white mb-3 text-shadow animate-floating-text">
            Breal Wallet
          </h1>

          <p
            style={{
              animationDelay: '1s',
            }}
            className="opacity-0 text-white text-lg overflow-hidden border-r-2 border-solid whitespace-nowrap animate-writer"
          >
            Your gateway to the decentralized future
          </p>

          <div className="w-16 h-1 bg-gradient-to-r from-white to-gray-300 mx-auto mt-4 rounded-full"></div>
        </div>
      </div>
      <div className="space-y-4 ">
        <Button
          variant="secondary"
          size="lg"
          onClick={onShowCreateWalletModal}
          className="text-sm transition-all duration-300"
        >
          <div className="flex items-center justify-center space-x-3">
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"></path>
            </svg>
            <span>Create New Wallet</span>
          </div>
        </Button>

        <Button
          variant="glass"
          size="lg"
          onClick={onShowImportWalletModal}
          className="rounded-2xl text-white font-semibold transition-all duration-300 text-sm"
        >
          <div className="flex items-center justify-center space-x-3">
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
              <path d="M14,2H6A2,2 0 0,0 4,4V20A2,2 0 0,0 6,22H18A2,2 0 0,0 20,20V8L14,2M18,20H6V4H13V9H18V20Z"></path>
            </svg>
            <span>Import Existing Wallet</span>
          </div>
        </Button>
      </div>

      <div className="animate-breathing-grow rounded-xl p-4 bg-[rgba(255, 255, 255, 0.25)] backdrop-blur-[20px] border-[1px] border-solid border-[rgba(255,255,255,0.3)] mt-8 shadow-[0px_25px_45px_rgba(0,0,0,0.1)]">
        <div className="flex items-center space-x-3">
          <div className="w-8 h-8 rounded-full bg-gradient-to-r from-green-400 to-blue-400 flex items-center justify-center">
            <svg
              className="w-4 h-4 text-white"
              fill="currentColor"
              viewBox="0 0 24 24"
            >
              <path d="M12,1L3,5V11C3,16.55 6.84,21.74 12,23C17.16,21.74 21,16.55 21,11V5L12,1M10,17L6,13L7.41,11.59L10,14.17L16.59,7.58L18,9L10,17Z"></path>
            </svg>
          </div>
          <div className="text-left">
            <p className="text-white font-medium text-sm animate-floating-text">
              Bank-Grade Security
            </p>
            <p className="text-white opacity-70 text-xs">
              End-to-end encrypted • Self-custody
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default memo(Welcome);
