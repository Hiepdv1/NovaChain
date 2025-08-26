import Header from './components/header';
import Footer from './components/footer';

const MainLayout = ({ children }: { children: React.ReactNode }) => {
  return (
    <div className="flex flex-col min-h-screen bg-gradient-to-br from-blue-200/70 via-indigo-200/70 to-purple-200/70 dark:from-gray-800/50 dark:via-gray-800/50 dark:to-gray-800/50 transition-colors duration-300">
      <Header />
      <main className="flex-1">{children}</main>
      <Footer />
    </div>
  );
};

export default MainLayout;
