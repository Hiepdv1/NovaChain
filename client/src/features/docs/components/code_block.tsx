import { CheckCircle2, Copy } from 'lucide-react';
import { memo } from 'react';

export interface CodeBlockProps {
  code: string;
  id: string;
  copiedCommand?: string;
  onCopy?: (code: string, id: string) => void;
}

const CodeBlock = ({ code, id, copiedCommand, onCopy }: CodeBlockProps) => (
  <div className="relative group">
    <div className="absolute right-2 top-2 z-10">
      <button
        onClick={() => onCopy?.(code, id)}
        className="px-3 py-1.5 rounded-lg bg-gray-700 hover:bg-gray-600 text-gray-300 text-sm flex items-center gap-1.5 transition-all opacity-0 group-hover:opacity-100"
      >
        {copiedCommand === id ? (
          <>
            <CheckCircle2 className="w-4 h-4 text-green-400" />
            <span className="text-green-400">Copied!</span>
          </>
        ) : (
          <>
            <Copy className="w-4 h-4" />
            <span>Copy</span>
          </>
        )}
      </button>
    </div>
    <pre className="bg-gray-900 dark:bg-gray-950 rounded-xl p-4 overflow-x-auto border border-gray-800 dark:border-gray-700">
      <code className="text-sm text-gray-300 font-mono">{code}</code>
    </pre>
  </div>
);

export default memo(CodeBlock);
