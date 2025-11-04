import { Tooltip } from '@/components/ui/tooltip';
import { FeeList } from '@/shared/constants/transaction';
import { TooltipContent, TooltipTrigger } from '@radix-ui/react-tooltip';
import { ChangeEvent } from 'react';

interface TransactionFeeProps {
  onChecked: (e: ChangeEvent<HTMLInputElement>) => void;
}

const TransactionFee = ({ onChecked }: TransactionFeeProps) => {
  return (
    <div className="space-y-4">
      <div className="flex items-center space-x-2">
        <label className="block text-xs font-medium text-gray-900 dark:text-white">
          Transaction Fee & Priority
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
                d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-8-3a1 1 0 00-.867.5 1 1 0 11-1.731-1A3 3 0 0113 8a3.001 3.001 0 01-2 2.83V11a1 1 0 11-2 0v-1a1 1 0 011-1 1 1 0 100-2zm0 8a1 1 0 100-2 1 1 0 000 2z"
                clipRule="evenodd"
              ></path>
            </svg>
          </TooltipTrigger>
          <TooltipContent
            align="center"
            className="bg-gray-600 max-w-md dark:bg-white/95 p-4 text-xs"
          >
            <strong className="text-white dark:text-black">
              PoW Fee Structure
            </strong>
            <ul className="list-disc list-inside mt-2 space-y-1 text-gray-300 dark:text-gray-950">
              <li>
                <strong>Slow</strong>: Minimum fee, lowest priority in the queue
              </li>
              <li>
                <strong>Standard</strong>: Balanced fee with medium priority
              </li>
              <li>
                <strong>Fast</strong>: Higher fee, prioritized for earlier
                processing
              </li>
              <li>
                Fees must be greater than <strong>0</strong>. The actual mining
                time cannot be predicted — fees only affect how your transaction
                is prioritized in the backend queue.
              </li>
            </ul>
          </TooltipContent>
        </Tooltip>
      </div>

      <div className="space-x-3">
        <div className="grid grid-cols-1 gap-3">
          {FeeList.map((fee) => {
            return (
              <label
                key={fee.id}
                className="before:content before:absolute before:top-0 before:-left-full before:w-full before:h-full before:bg-[linear-gradient(90deg,_transparent,_rgba(255,255,255,0.2),_transparent)] before:transition-all before:duration-500 hover:before:left-full relative cursor-pointer overflow-hidden"
                htmlFor={fee.value}
              >
                <input
                  value={fee.value}
                  id={fee.value}
                  className="sr-only peer"
                  type="radio"
                  name="fee"
                  defaultChecked={fee.checked}
                  onChange={onChecked}
                />
                <div className="p-4 bg-white/50 dark:bg-gray-800/50 border border-gray-200/50 dark:border-gray-800/50 rounded-2xl peer-checked:!border-primary-500 peer-checked:bg-primary-50/50 peer-checked:dark:bg-primary-900/20 transition-all duration-200 hover:bg-white/70 dark:hover:bg-gray-800/70">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-3">
                      <div
                        className={`w-3 h-3 ${fee.color} rounded-full`}
                      ></div>
                      <div>
                        <div className="text-sm font-medium text-gray-900 dark:text-white">
                          {fee.title}
                        </div>
                        <div className="text-xs text-gray-500 dark:text-gray-400">
                          {fee.des}
                        </div>
                      </div>
                    </div>

                    <div className="text-right">
                      <div className="text-xs font-semibold text-gray-900 dark:text-white">
                        {fee.fee} CCC
                      </div>
                    </div>
                  </div>
                </div>
              </label>
            );
          })}
        </div>
      </div>
    </div>
  );
};

export default TransactionFee;
