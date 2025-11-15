'use client';

import { Book, Cpu, Download, Menu, Shield, Terminal, Zap } from 'lucide-react';
import { memo, useState } from 'react';
import { DocContent, DocSection, DowloadProps } from '../types/docs';
import CodeBlock, { CodeBlockProps } from './code_block';
import GettingStartedPage from './getting-started';
import InstallationPage from './installation-page';
import CLIReferencePage from './cli-reference';
import WalletPage from './wallet-page';
import MiningPage from './mining-page';
import ApiPage from './api-page';
import Sidebar from './sidebar';

interface NovaChainProps {
  SizeFile: number;
}

const NovaChain = ({ SizeFile }: NovaChainProps) => {
  const [sidebarOpen, setSidebarOpen] = useState<boolean>(false);
  const [activeSection, setActiveSection] =
    useState<DocContent>('getting-started');
  const [copiedCommand, setCopiedCommand] = useState<string>('');

  const downloads: DowloadProps[] = [
    // { platform: 'Linux', version: 'v1.0.0', size: '45 MB', icon: '🐧' },
    // { platform: 'macOS', version: 'v1.0.0', size: '42 MB', icon: '🍎' },
    {
      platform: 'Windows',
      version: 'v1.0.0',
      size: `${SizeFile} MB`,
      icon: '🪟',
      name: 'novachain.rar',
    },
  ];

  const docSections: DocSection[] = [
    { id: 'getting-started', label: 'Getting Started', icon: Zap },
    { id: 'installation', label: 'Installation', icon: Download },
    { id: 'cli-reference', label: 'CLI Reference', icon: Terminal },
    { id: 'wallet', label: 'Wallet Management', icon: Shield },
    { id: 'mining', label: 'Mining', icon: Cpu },
    { id: 'api', label: 'JSON-RPC API', icon: Book },
  ];

  const copyToClipboard = (text: string, id: string) => {
    navigator.clipboard.writeText(text);
    setCopiedCommand(id);
    setTimeout(() => setCopiedCommand(''), 2000);
  };

  const CodeBlockWithCopy = (props: CodeBlockProps) => (
    <CodeBlock
      {...props}
      copiedCommand={copiedCommand}
      onCopy={copyToClipboard}
    />
  );

  const renderContent = () => {
    const commonProps = { CodeBlock: CodeBlockWithCopy };

    const pages = {
      'getting-started': (
        <GettingStartedPage {...commonProps} onNavigate={setActiveSection} />
      ),
      installation: <InstallationPage {...commonProps} downloads={downloads} />,
      'cli-reference': <CLIReferencePage {...commonProps} />,
      wallet: <WalletPage {...commonProps} />,
      mining: <MiningPage {...commonProps} />,
      api: <ApiPage {...commonProps} />,
    };

    return pages[activeSection] || pages['getting-started'];
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 transition-colors duration-300">
      <div className="flex">
        {/* Sidebar */}
        <Sidebar
          docSections={docSections}
          activeSection={activeSection}
          onNavigate={setActiveSection}
          sidebarOpen={sidebarOpen}
          onClose={() => setSidebarOpen(false)}
        />

        {/* Overlay for mobile */}
        {sidebarOpen && (
          <div
            className="fixed inset-0 bg-black/50 z-40 lg:hidden"
            onClick={() => setSidebarOpen(false)}
          />
        )}

        {/* Main Content */}
        <main className="flex-1 min-w-0">
          {/* Mobile Header */}
          <div className="lg:hidden sticky top-0 z-30 bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 px-4 py-3">
            <div className="flex items-center justify-between">
              <button
                onClick={() => setSidebarOpen(true)}
                className="p-2 rounded-lg bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 dark:hover:bg-gray-600 transition-all"
              >
                <Menu className="w-5 h-5 text-gray-700 dark:text-gray-300" />
              </button>
              <div className="flex items-center gap-2">
                <button className="flex items-center gap-2 px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-lg font-medium transition-all text-sm">
                  <Download className="w-4 h-4" />
                  Download
                </button>
              </div>
            </div>
          </div>

          {/* Content Area */}
          <div className="p-6 lg:p-8">
            <div className="max-w-4xl">
              {renderContent()}

              {/* Footer */}
              <div className="mt-12 pt-8 border-t border-gray-200 dark:border-gray-700">
                <div className="flex flex-col sm:flex-row justify-between items-center gap-4 text-sm text-gray-600 dark:text-gray-400">
                  <div>© 2024 NovaChain. Open source project.</div>
                  <div className="flex gap-4">
                    <a
                      href="#"
                      className="hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                    >
                      GitHub
                    </a>
                    <a
                      href="#"
                      className="hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                    >
                      License
                    </a>
                    <a
                      href="#"
                      className="hover:text-blue-600 dark:hover:text-blue-400 transition-colors"
                    >
                      Privacy
                    </a>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </main>
      </div>
    </div>
  );
};

export default memo(NovaChain);
