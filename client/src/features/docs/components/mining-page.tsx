import { ComponentType, memo } from 'react';
import { CodeBlockProps } from './code_block';

interface MiningPageProps {
  CodeBlock: ComponentType<CodeBlockProps>;
}

const MiningPage = ({ CodeBlock }: MiningPageProps) => (
  <div className="space-y-8">
    <div>
      <h1 className="text-4xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        Mining Guide
      </h1>
      <p className="text-lg text-gray-600 dark:text-gray-400">
        Learn how to start mining blocks on NovaChain.
      </p>
    </div>

    <div>
      <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        Starting a Mining Node
      </h2>
      <p className="text-gray-600 dark:text-gray-400 mb-4">
        To start mining, you need a wallet address to receive mining rewards:
      </p>
      <CodeBlock
        code={`# First, create a wallet if you don't have one
novachain wallet new

# Start mining node
novachain startNode \\
  --Port 3000 \\
  --InstanceId 1001 \\
  --Miner true \\
  --Address <your_wallet_address>`}
        id="start-mining"
      />
    </div>

    <div>
      <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        Mining with RPC Enabled
      </h2>
      <p className="text-gray-600 dark:text-gray-400 mb-4">
        Enable JSON-RPC for monitoring and management:
      </p>
      <CodeBlock
        code={`novachain startNode \\
  --Port 3000 \\
  --InstanceId 1001 \\
  --Miner true \\
  --Address <your_wallet_address> \\
  --RPC \\
  --RPC-Port 9000`}
        id="mining-rpc"
      />
    </div>

    <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
      <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        Mining Rewards
      </h3>
      <div className="space-y-3 text-gray-600 dark:text-gray-400">
        <p>
          Mining rewards are automatically sent to your specified wallet address
          when you successfully mine a block.
        </p>
        <div className="grid grid-cols-2 gap-4 mt-4">
          <div className="bg-gray-50 dark:bg-gray-900 p-4 rounded-lg">
            <div className="text-sm text-gray-500 dark:text-gray-400">
              Block Reward
            </div>
            <div className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              50 COIN
            </div>
          </div>
          <div className="bg-gray-50 dark:bg-gray-900 p-4 rounded-lg">
            <div className="text-sm text-gray-500 dark:text-gray-400">
              Avg Block Time
            </div>
            <div className="text-2xl font-bold text-gray-900 dark:text-gray-100">
              ~120s
            </div>
          </div>
        </div>
      </div>
    </div>

    <div className="bg-green-50 dark:bg-green-900/20 border-l-4 border-green-500 p-6 rounded-r-lg">
      <h3 className="text-lg font-semibold text-green-900 dark:text-green-100 mb-2">
        💡 Pro Tips
      </h3>
      <ul className="list-disc list-inside text-green-800 dark:text-green-200 space-y-2">
        <li>Run multiple mining nodes for better chances</li>
        <li>Monitor your node performance with RPC</li>
        <li>Keep your node synced with the network</li>
        <li>Use SSD storage for better performance</li>
      </ul>
    </div>
  </div>
);

export default memo(MiningPage);
