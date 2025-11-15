import { ComponentType, memo } from 'react';
import { CodeBlockProps } from './code_block';
import { ChevronRight } from 'lucide-react';
import { DocContent } from '../types/docs';

interface GettingStartedProps {
  onNavigate: (path: DocContent) => void;
  CodeBlock: ComponentType<CodeBlockProps>;
}

const GettingStartedPage = ({ onNavigate, CodeBlock }: GettingStartedProps) => (
  <div className="space-y-8">
    <div>
      <h1 className="text-4xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        Getting Started
      </h1>
      <p className="text-lg text-gray-600 dark:text-gray-400 leading-relaxed">
        Welcome to NovaChain! This guide will help you get up and running with
        your own private blockchain network.
      </p>
    </div>

    <div className="bg-blue-50 dark:bg-blue-900/20 border-l-4 border-blue-500 p-6 rounded-r-lg">
      <h3 className="text-lg font-semibold text-blue-900 dark:text-blue-100 mb-2">
        What is NovaChain?
      </h3>
      <p className="text-blue-800 dark:text-blue-200">
        NovaChain is a private Proof-of-Work (PoW) blockchain implementation
        designed for enterprise solutions with enhanced security, scalability,
        and performance.
      </p>
    </div>

    <div>
      <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        Quick Start
      </h2>
      <div className="space-y-4">
        <div className="flex items-start gap-4">
          <div className="flex-shrink-0 w-8 h-8 rounded-full bg-blue-600 text-white flex items-center justify-center font-bold">
            1
          </div>
          <div className="flex-1">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2">
              Download NovaChain
            </h3>
            <p className="text-gray-600 dark:text-gray-400 mb-3">
              Download the latest version for your platform from the downloads
              section.
            </p>
            <button
              onClick={() => onNavigate('installation')}
              className="cursor-pointer text-blue-600 dark:text-blue-400 hover:text-blue-700 dark:hover:text-blue-300 font-medium inline-flex items-center gap-1"
            >
              View Installation Guide
              <ChevronRight className="w-4 h-4" />
            </button>
          </div>
        </div>

        <div className="flex items-start gap-4">
          <div className="flex-shrink-0 w-8 h-8 rounded-full bg-blue-600 text-white flex items-center justify-center font-bold">
            2
          </div>
          <div className="flex-1">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2">
              Initialize Blockchain
            </h3>
            <p className="text-gray-600 dark:text-gray-400 mb-3">
              Create your first blockchain instance with the genesis block.
            </p>
            <CodeBlock code="novachain init --InstanceId 1001" id="init-cmd" />
          </div>
        </div>

        <div className="flex items-start gap-4">
          <div className="flex-shrink-0 w-8 h-8 rounded-full bg-blue-600 text-white flex items-center justify-center font-bold">
            3
          </div>
          <div className="flex-1">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2">
              Create a Wallet
            </h3>
            <p className="text-gray-600 dark:text-gray-400 mb-3">
              Generate a new wallet to send and receive transactions.
            </p>
            <CodeBlock code="novachain wallet new" id="wallet-cmd" />
          </div>
        </div>

        <div className="flex items-start gap-4">
          <div className="flex-shrink-0 w-8 h-8 rounded-full bg-blue-600 text-white flex items-center justify-center font-bold">
            4
          </div>
          <div className="flex-1">
            <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-2">
              Start a Node
            </h3>
            <p className="text-gray-600 dark:text-gray-400 mb-3">
              Launch a full node or mining node.
            </p>
            <CodeBlock
              code="novachain startNode --Port 3000 --InstanceId 1001"
              id="node-cmd"
            />
          </div>
        </div>
      </div>
    </div>

    <div className="bg-yellow-50 dark:bg-yellow-900/20 border-l-4 border-yellow-500 p-6 rounded-r-lg">
      <h3 className="text-lg font-semibold text-yellow-900 dark:text-yellow-100 mb-2 flex items-center gap-2">
        <span>⚠️</span>
        Prerequisites
      </h3>
      <ul className="list-disc list-inside text-yellow-800 dark:text-yellow-200 space-y-1">
        <li>Go 1.19 or higher (for building from source)</li>
        <li>At least 4GB RAM</li>
        <li>10GB free disk space</li>
      </ul>
    </div>
  </div>
);

export default memo(GettingStartedPage);
