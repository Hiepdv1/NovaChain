export interface ReAuthModal {
  title: string;
  des: string;
  notice?: React.ReactNode;
  wallet: {
    address: string;
  };
}
