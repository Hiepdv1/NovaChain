import ActivityComponent from '../components/Activity';
import OverviewComponent from '../components/overview';

const HomePage = () => {
  return (
    <section className="glass-card overflow-hidden dark:bg-primary-dark/5 dark:border-secondary-dark rounded-2xl">
      <div className="animate-dashboard-fade-in p-6">
        <OverviewComponent />
        <ActivityComponent />
      </div>
    </section>
  );
};

export default HomePage;
