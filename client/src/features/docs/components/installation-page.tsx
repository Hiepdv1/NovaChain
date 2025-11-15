import { ComponentType, memo } from 'react';
import { CodeBlockProps } from './code_block';
import { DowloadProps } from '../types/docs';
import Download_card from './download_card';

interface InstallationProps {
  CodeBlock: ComponentType<CodeBlockProps>;
  downloads: DowloadProps[];
}

const Installation = ({ downloads, CodeBlock }: InstallationProps) => {
  const onDownload = (file: string) => {
    window.location.href = `${process.env.NEXT_PUBLIC_API_URL}download/${file}`;
  };

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-4xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          Installation
        </h1>
        <p className="text-lg text-gray-600 dark:text-gray-400">
          Install NovaChain on your preferred platform.
        </p>
      </div>

      <div>
        <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-6">
          Download Binary
        </h2>
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
          {downloads.map((dl, idx) => (
            <Download_card
              key={idx}
              {...dl}
              onDownload={() => onDownload(dl.name)}
            />
          ))}
        </div>
      </div>

      <div>
        <h2 className="text-2xl font-bold text-gray-900 dark:text-gray-100 mb-4">
          Installation Steps
        </h2>

        <div className="space-y-6">
          <div>
            <h3 className="text-xl font-semibold text-gray-900 dark:text-gray-100 mb-3">
              Windows
            </h3>
            <CodeBlock
              code={`# Extract the ZIP file
# Add novachain.exe to your PATH

# Verify installation (in PowerShell)
novachain --help`}
              id="windows-install"
            />
          </div>
        </div>
      </div>
    </div>
  );
};

export default memo(Installation);
