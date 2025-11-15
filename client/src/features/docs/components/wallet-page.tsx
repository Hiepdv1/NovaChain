import { ComponentType } from 'react';
import { CodeBlockProps } from './code_block';

interface WalletPageProps {
  CodeBlock: ComponentType<CodeBlockProps>;
}

const WalletPage = ({ CodeBlock }: WalletPageProps) => (
  <div className="space-y-8">
    <div>
      <h1 className="text-4xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        Wallet Management
      </h1>
      <p className="text-lg text-gray-600 dark:text-gray-400">
        Learn how to create and manage wallets in NovaChain.
      </p>
    </div>

    <div>
      <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        Creating a Wallet
      </h2>
      <p className="text-gray-600 dark:text-gray-400 mb-4">
        Generate a new wallet with public/private key pair:
      </p>
      <CodeBlock code="novachain wallet new" id="create-wallet" />
      <div className="mt-4 bg-blue-50 dark:bg-blue-900/20 border-l-4 border-blue-500 p-4 rounded-r-lg">
        <p className="text-sm text-blue-800 dark:text-blue-200">
          <strong>Important:</strong> Save your private key securely! It cannot
          be recovered if lost.
        </p>
      </div>
    </div>

    <div>
      <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        Listing Wallets
      </h2>
      <p className="text-gray-600 dark:text-gray-400 mb-4">
        View all wallet addresses in your system:
      </p>
      <CodeBlock code="novachain wallet list" id="list-wallet" />
    </div>

    <div>
      <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        Checking Balance
      </h2>
      <p className="text-gray-600 dark:text-gray-400 mb-4">
        Check the balance of a specific wallet:
      </p>
      <CodeBlock
        code={`novachain wallet balance \\
  --Address 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb \\
  --InstanceId 1001`}
        id="check-balance"
      />
    </div>

    <div className="bg-yellow-50 dark:bg-yellow-900/20 border-l-4 border-yellow-500 p-6 rounded-r-lg">
      <h3 className="text-lg font-semibold text-yellow-900 dark:text-yellow-100 mb-2">
        Security Best Practices
      </h3>
      <ul className="list-disc list-inside text-yellow-800 dark:text-yellow-200 space-y-2">
        <li>Never share your private key with anyone</li>
        <li>Store your private key in a secure location</li>
        <li>Use hardware wallets for large amounts</li>
        <li>Enable backup and recovery mechanisms</li>
      </ul>
    </div>
  </div>
);

export default WalletPage;
