'use client';

import { CheckCircle2, Copy } from 'lucide-react';
import { useState } from 'react';

interface CopyButton {
  text: string;
  label: string;
}

const CopyButton = ({ text, label }: CopyButton) => {
  const [copied, setCopied] = useState(false);

  const handleCopy = () => {
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <button
      onClick={handleCopy}
      className="cursor-pointer inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 transition-all text-sm"
      title={`Copy ${label}`}
    >
      {copied ? (
        <>
          <CheckCircle2 className="w-4 h-4 text-green-600 dark:text-green-400" />
          <span className="text-green-600 dark:text-green-400 font-medium">
            Copied!
          </span>
        </>
      ) : (
        <>
          <Copy className="w-4 h-4 text-gray-600 dark:text-gray-400" />
          <span className="text-gray-700 dark:text-gray-300">Copy</span>
        </>
      )}
    </button>
  );
};

export default CopyButton;
