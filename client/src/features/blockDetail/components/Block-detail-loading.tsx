const BlockDetailSkeletonLoader = () => (
  <div className="animate-pulse space-y-6">
    <div className={`h-32 rounded-xl bg-slate-200 dark:bg-slate-700`}></div>
    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
      {[...Array(8)].map((_, i) => (
        <div
          key={i}
          className={`h-24 rounded-xl bg-slate-200 dark:bg-slate-700`}
        ></div>
      ))}
    </div>
    <div className="space-y-4">
      {[...Array(5)].map((_, i) => (
        <div
          key={i}
          className={`h-32 rounded-xl dark:bg-slate-700 bg-slate-200`}
        ></div>
      ))}
    </div>
  </div>
);

export default BlockDetailSkeletonLoader;
