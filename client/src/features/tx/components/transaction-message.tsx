import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { ChangeEvent } from 'react';

interface TransactionMessageProps {
  onInputMessage?: (e: ChangeEvent<HTMLTextAreaElement>) => void;
}

const TransactionMessage = ({ onInputMessage }: TransactionMessageProps) => {
  return (
    <div className="space-y-2">
      <div className="flex items-center space-x-2">
        <label
          htmlFor="message"
          className="block text-xs font-medium text-gray-900 dark:text-white"
        >
          Message (Optional)
        </label>
        <Tooltip>
          <TooltipTrigger asChild>
            <svg
              className="w-4 h-4 text-gray-400 cursor-help"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-3a1 1 0 00-.867.5 1 1 0 11-1.731-1A3 3 0 0113 8a3 3 0 01-2 2.83V11a1 1 0 11-2 0v-1a1 1 0 011-1 1 1 0 100-2zm0 8a1 1 0 100-2 1 1 0 000 2z"
                clipRule="evenodd"
              ></path>
            </svg>
          </TooltipTrigger>
          <TooltipContent
            align="center"
            className="bg-gray-600 dark:bg-white/95 p-4 rounded-lg shadow-lg"
          >
            <strong>Private Transaction Message</strong>
            <ul className="list-disc list-inside mt-2 space-y-1 text-gray-300 dark:text-gray-950">
              <li>
                Add a private note <strong>(max 1000 characters)</strong>
              </li>
              <li>End-to-end encrypted between sender and recipient</li>
              <li>Stored securely off-chain</li>
              <li>Not visible on blockchain explorers</li>
              <li>For your personal reference only</li>
            </ul>
          </TooltipContent>
        </Tooltip>
      </div>

      <div
        className="pl-4 pr-2 py-3 bg-white/70 dark:bg-gray-800/70 
             border border-gray-200/50 dark:border-gray-600/50 
             rounded-2xl backdrop-blur-sm 
             transition-all duration-200 
             focus-within:ring-2 focus-within:ring-primary-500/50 
             focus-within:border-primary-500/50"
      >
        <textarea
          onChange={onInputMessage}
          className="text-sm w-full text-gray-900 dark:text-white 
               placeholder:text-gray-500 dark:placeholder:text-gray-400 
               focus:outline-none resize-none"
          name="message"
          id="message"
          rows={5}
          maxLength={1000}
          placeholder="Add a note to this transaction (optional)"
        ></textarea>
      </div>

      <div className="flex items-center justify-between text-xs text-gray-500 dark:text-gray-400">
        <span>Encrypted private note — only visible to the recipient</span>
        <span>0/1000</span>
      </div>
    </div>
  );
};

export default TransactionMessage;
