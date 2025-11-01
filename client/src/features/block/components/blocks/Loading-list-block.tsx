const ListBlockLoadingSkeleton = () => {
  return (
    <div className="animate-pulse space-y-4">
      {[...Array(9)].map((_, i) => (
        <div
          key={i}
          className="bg-gray-100 dark:bg-gray-800 rounded-xl h-48 transition-colors duration-300"
        />
      ))}
    </div>
  );
};

export default ListBlockLoadingSkeleton;
