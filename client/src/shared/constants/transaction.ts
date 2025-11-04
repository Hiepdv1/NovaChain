type FeeoptionType = 'slow' | 'standard' | 'fast';

export interface FeeOption {
  id: string;
  title: string;
  des: string;
  fee: number;
  value: FeeoptionType;
  color: string;
  priority: number;
  checked: boolean;
}

export const FeeList: FeeOption[] = [
  {
    id: '1',
    title: 'Slow (Economy)',
    des: 'Lower fee • Waits for less busy blocks',
    fee: 0.1,
    value: 'slow',
    color: 'bg-yellow-500',
    priority: 1,
    checked: false,
  },
  {
    id: '2',
    title: 'Standard (Recommended)',
    des: 'Balanced fee • Normal mining priority',
    fee: 0.5,
    value: 'standard',
    color: 'bg-blue-500',
    priority: 2,
    checked: true,
  },
  {
    id: '3',
    title: 'Fast (Priority)',
    des: 'Higher fee • Miners prioritize your transaction',
    fee: 0.9,
    value: 'fast',
    color: 'bg-green-500',
    priority: 3,
    checked: false,
  },
];
