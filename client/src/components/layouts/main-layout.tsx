import Header from './components/header';
import Footer from './components/footer';
import NavBar from './components/navbar';

const MainLayout = ({ children }: { children: React.ReactNode }) => {
  return (
    <div className="flex flex-col min-h-screen bg-gradient-to-br from-blue-200/70 via-indigo-200/70 to-purple-200/70 dark:from-gray-800/50 dark:via-gray-800/50 dark:to-gray-800/50 transition-colors duration-300">
      <Header />
      <main className="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 flex-1">
        <NavBar />
        {children}
      </main>
      <Footer />
    </div>
  );
};

export default MainLayout;
