import { ComponentType, memo } from 'react';
import { CodeBlockProps } from './code_block';

interface CLIReferenceProps {
  CodeBlock: ComponentType<CodeBlockProps>;
}

const CLIReferencePage = ({ CodeBlock }: CLIReferenceProps) => (
  <div className="space-y-8">
    <div>
      <h1 className="text-4xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        CLI Reference
      </h1>
      <p className="text-lg text-gray-600 dark:text-gray-400">
        Complete command-line interface reference for NovaChain.
      </p>
    </div>

    <div className="space-y-6">
      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
        <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          novachain init
        </h3>
        <p className="text-gray-600 dark:text-gray-400 mb-4">
          Initialize a new blockchain instance with a genesis block.
        </p>
        <CodeBlock code="novachain init --InstanceId <id>" id="init-ref" />
        <div className="mt-4">
          <h4 className="font-semibold text-gray-900 dark:text-gray-100 mb-2">
            Flags:
          </h4>
          <ul className="space-y-2 text-sm text-gray-600 dark:text-gray-400">
            <li>
              <code className="px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded">
                --InstanceId
              </code>{' '}
              - Unique blockchain instance ID (auto-generated if omitted)
            </li>
          </ul>
        </div>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
        <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          novachain startNode
        </h3>
        <p className="text-gray-600 dark:text-gray-400 mb-4">
          Start a full node or mining node.
        </p>
        <CodeBlock
          code="novachain startNode --Port 3000 --InstanceId 1001 [options]"
          id="node-ref"
        />
        <div className="mt-4">
          <h4 className="font-semibold text-gray-900 dark:text-gray-100 mb-2">
            Required Flags:
          </h4>
          <ul className="space-y-2 text-sm text-gray-600 dark:text-gray-400 mb-4">
            <li>
              <code className="px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded">
                --Port
              </code>{' '}
              - Node listening port
            </li>
            <li>
              <code className="px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded">
                --InstanceId
              </code>{' '}
              - Blockchain instance ID
            </li>
          </ul>
          <h4 className="font-semibold text-gray-900 dark:text-gray-100 mb-2">
            Optional Flags:
          </h4>
          <ul className="space-y-2 text-sm text-gray-600 dark:text-gray-400">
            <li>
              <code className="px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded">
                --Miner
              </code>{' '}
              - Enable mining mode (default: false)
            </li>
            <li>
              <code className="px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded">
                --Address
              </code>{' '}
              - Miner wallet address (required if Miner=true)
            </li>
            <li>
              <code className="px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded">
                --Fullnode
              </code>{' '}
              - Run as full node (default: true)
            </li>
            <li>
              <code className="px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded">
                --SeedPeer
              </code>{' '}
              - Enable seed peer discovery
            </li>
            <li>
              <code className="px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded">
                --RPC
              </code>{' '}
              - Enable JSON-RPC server
            </li>
            <li>
              <code className="px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded">
                --RPC-Port
              </code>{' '}
              - RPC server port (default: 9000)
            </li>
            <li>
              <code className="px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded">
                --RPC-Addr
              </code>{' '}
              - RPC server address (default: 127.0.0.1)
            </li>
            <li>
              <code className="px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded">
                --RPC-Mode
              </code>{' '}
              - RPC mode: http, tcp, both
            </li>
          </ul>
        </div>
      </div>

      <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
        <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          novachain wallet
        </h3>
        <p className="text-gray-600 dark:text-gray-400 mb-4">
          Wallet management commands.
        </p>

        <div className="space-y-4">
          <div>
            <h4 className="font-semibold text-gray-900 dark:text-gray-100 mb-2">
              Create new wallet:
            </h4>
            <CodeBlock code="novachain wallet new" id="wallet-new" />
          </div>

          <div>
            <h4 className="font-semibold text-gray-900 dark:text-gray-100 mb-2">
              List all wallets:
            </h4>
            <CodeBlock code="novachain wallet list" id="wallet-list" />
          </div>

          <div>
            <h4 className="font-semibold text-gray-900 dark:text-gray-100 mb-2">
              Check balance:
            </h4>
            <CodeBlock
              code="novachain wallet balance --Address <wallet_address> --InstanceId <id>"
              id="wallet-balance"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
);

export default memo(CLIReferencePage);
