import { usePathname, useSearchParams } from 'next/navigation';

const useCurrentUrl = () => {
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const query = searchParams.toString();

  const fullPath = query ? `${pathname}?${query}` : pathname;
  return fullPath;
};

export default useCurrentUrl;
