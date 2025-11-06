const RankBadge = ({ rank }: { rank: number }) => {
  const getRankStyle = () => {
    if (rank === 1)
      return 'bg-gradient-to-r from-yellow-400 to-yellow-500 text-yellow-900';
    if (rank === 2)
      return 'bg-gradient-to-r from-gray-300 to-gray-400 text-gray-800';
    if (rank === 3)
      return 'bg-gradient-to-r from-orange-400 to-orange-500 text-orange-900';
    return 'bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300';
  };

  return (
    <div
      className={`flex items-center justify-center w-12 h-12 rounded-xl font-bold text-lg shadow-md ${getRankStyle()}`}
    >
      #{rank}
    </div>
  );
};

export default RankBadge;
