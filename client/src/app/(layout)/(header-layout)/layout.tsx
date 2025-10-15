import HeaderLayout from '@/components/layouts/header-layout';

const Layout = ({ children }: { children: React.ReactNode }) => {
  return <HeaderLayout>{children}</HeaderLayout>;
};

export default Layout;
