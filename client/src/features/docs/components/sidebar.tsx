import { memo } from 'react';
import { DocContent, DocSection } from '../types/docs';

interface SidebarProps {
  docSections: DocSection[];
  activeSection: string;
  onNavigate: (content: DocContent) => void;
  sidebarOpen: boolean;
  onClose: () => void;
}

const Sidebar = ({
  docSections,
  activeSection,
  onNavigate,
  sidebarOpen,
  onClose,
}: SidebarProps) => (
  <aside
    className={`
    fixed lg:static inset-y-0 left-0 z-50 w-72
    bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700
    transform transition-transform duration-300 ease-in-out
    ${sidebarOpen ? 'translate-x-0' : '-translate-x-full lg:translate-x-0'}
  `}
  >
    <div className="h-full flex flex-col">
      <nav className="flex-1 overflow-y-auto p-4">
        <div className="space-y-1">
          {docSections.map((section) => {
            const Icon = section.icon;
            return (
              <button
                key={section.id}
                onClick={() => {
                  onNavigate(section.id);
                  onClose();
                }}
                className={`cursor-pointer w-full flex items-center gap-3 px-4 py-3 rounded-lg text-left transition-all ${
                  activeSection === section.id
                    ? 'bg-blue-600 text-white shadow-lg'
                    : 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700'
                }`}
              >
                <Icon className="w-5 h-5" />
                <span className="font-medium">{section.label}</span>
              </button>
            );
          })}
        </div>
      </nav>
    </div>
  </aside>
);

export default memo(Sidebar);
