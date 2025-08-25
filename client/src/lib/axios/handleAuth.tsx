import { useModalStore } from '@/stores/modal-store';
import { GetWalletPool } from '../db/wallet.index';

export const handleUnauthorized = async () => {
  const walletPool = await GetWalletPool();

  if (walletPool.length > 0) {
    useModalStore.getState().actions.openModal({
      type: 'reauth',
      props: {
        des: 'Please enter your password to verify your identity and continue with this action.',
        title: 'Confirm Your Identity',
        notice: (
          <div className="mt-6 p-4 bg-gray-50 rounded-xl shadow-white">
            <div className="flex items-center text-gray-600 text-xs">
              <svg
                className="w-4 h-4 mr-2"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path
                  fillRule="evenodd"
                  d="M5 9V7a5 5 0 0110 0v2a2 2 0 012 2v5a2 2 0 01-2 2H5a2 2 0 01-2-2v-5a2-2 0 012-2zm8-2v2H7V7a3 3 0 016 0z"
                  clipRule="evenodd"
                ></path>
              </svg>
              <span>
                Your information is protected with end-to-end encryption
              </span>
            </div>
          </div>
        ),
        wallet: {
          address: walletPool[0].address,
        },
      },
    });
  }
};
