import { ComponentType, memo } from 'react';
import { CodeBlockProps } from './code_block';

interface APIPageProps {
  CodeBlock: ComponentType<CodeBlockProps>;
}

const APIPage = ({ CodeBlock }: APIPageProps) => (
  <div className="space-y-8">
    <div>
      <h1 className="text-4xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        JSON-RPC API
      </h1>
      <p className="text-lg text-gray-600 dark:text-gray-400">
        Interact with NovaChain programmatically using JSON-RPC.
      </p>
    </div>

    <div>
      <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        Enabling RPC Server
      </h2>
      <CodeBlock
        code={`novachain startNode \\
  --Port 3000 \\
  --InstanceId 1001 \\
  --RPC \\
  --RPC-Port 9000 \\
  --RPC-Addr 127.0.0.1 \\
  --RPC-Mode http`}
        id="enable-rpc"
      />
    </div>

    <div>
      <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-4">
        API Endpoints
      </h2>
      <div className="space-y-4">
        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <h3 className="text-lg font-bold text-gray-900 dark:text-gray-100 mb-2">
            Get Block by Height
          </h3>
          <CodeBlock
            code={`curl -X POST http://localhost:9000 \\
  -H "Content-Type: application/json" \\
  -d '{
    "jsonrpc": "2.0",
    "method": "getBlock",
    "params": [12345],
    "id": 1
  }'`}
            id="api-getblock"
          />
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <h3 className="text-lg font-bold text-gray-900 dark:text-gray-100 mb-2">
            Get Transaction
          </h3>
          <CodeBlock
            code={`curl -X POST http://localhost:9000 \\
  -H "Content-Type: application/json" \\
  -d '{
    "jsonrpc": "2.0",
    "method": "getTransaction",
    "params": ["0x123..."],
    "id": 1
  }'`}
            id="api-gettx"
          />
        </div>

        <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border border-gray-200 dark:border-gray-700">
          <h3 className="text-lg font-bold text-gray-900 dark:text-gray-100 mb-2">
            Get Balance
          </h3>
          <CodeBlock
            code={`curl -X POST http://localhost:9000 \\
  -H "Content-Type: application/json" \\
  -d '{
    "jsonrpc": "2.0",
    "method": "getBalance",
    "params": ["0xAddress..."],
    "id": 1
  }'`}
            id="api-balance"
          />
        </div>
      </div>
    </div>
  </div>
);

export default memo(APIPage);
