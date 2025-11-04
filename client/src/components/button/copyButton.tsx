import { Check, CheckCircle2, Copy } from 'lucide-react';
import { MouseEvent, useState } from 'react';

interface CopyButtonProps {
  text: string;
  type?: 'text' | 'icon';
}

const CopyButton = ({ text, type = 'text' }: CopyButtonProps) => {
  const [copied, setCopied] = useState(false);

  if (type === 'icon') {
    const handleCopy = (e: MouseEvent<HTMLButtonElement>) => {
      e.stopPropagation();
      navigator.clipboard.writeText(text);
      setCopied(true);
      setTimeout(() => setCopied(false), 2000);
    };

    return (
      <button
        onClick={handleCopy}
        className="p-1.5 rounded-lg bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 transition-all hover:scale-110"
        title="Copy"
      >
        {copied ? (
          <CheckCircle2 className="w-4 h-4 text-green-600 dark:text-green-400" />
        ) : (
          <Copy className="w-4 h-4 text-gray-600 dark:text-gray-400" />
        )}
      </button>
    );
  }

  const handleCopy = () => {
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="flex items-center gap-2">
      <button
        onClick={handleCopy}
        className={`${
          copied ? '' : 'cursor-pointer'
        } text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors`}
        title="Copy"
      >
        {copied ? (
          <Check className="w-4 h-4 text-green-600" />
        ) : (
          <Copy className="w-4 h-4" />
        )}
      </button>
      {copied && (
        <span className="text-xs text-green-600 dark:text-green-400 whitespace-nowrap">
          Copied!
        </span>
      )}
    </div>
  );
};

export default CopyButton;
