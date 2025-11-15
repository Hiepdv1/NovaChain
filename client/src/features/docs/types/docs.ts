import { ComponentType } from 'react';

export interface DowloadProps {
  platform: string;
  version: string;
  size: string;
  icon: string;
  name: string;
}

export interface DocSection {
  id: DocContent;
  label: string;
  icon: ComponentType<{ className: string }>;
}

export type DocContent =
  | 'getting-started'
  | 'installation'
  | 'cli-reference'
  | 'wallet'
  | 'mining'
  | 'api';
