import { Download } from 'lucide-react';
import { memo } from 'react';

interface DownloadCardProps {
  platform: string;
  version: string;
  size: string;
  icon: string;
  onDownload: () => void;
}

const DownloadCard = ({
  platform,
  version,
  size,
  icon,
  onDownload,
}: DownloadCardProps) => (
  <div className="bg-white dark:bg-gray-800 rounded-xl p-6 border-2 border-gray-200 dark:border-gray-700 hover:border-blue-500 dark:hover:border-blue-500 transition-all">
    <div className="text-center">
      <div className="text-4xl mb-3">{icon}</div>
      <h3 className="text-xl font-bold text-gray-900 dark:text-gray-100 mb-2">
        {platform}
      </h3>
      <div className="text-sm text-gray-600 dark:text-gray-400 mb-4">
        <div>Version: {version}</div>
        <div>Size: {size}</div>
      </div>
    </div>
    <button
      onClick={onDownload}
      className="cursor-pointer w-full flex items-center justify-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-semibold transition-all"
    >
      <Download className="w-4 h-4" />
      Download
    </button>
  </div>
);

export default memo(DownloadCard);
