export const FormatHash = (hash: string) => {
  return `${hash.slice(0, 4)}...${hash.slice(-4)}`;
};

export const FormatTimestamp = (timestamp: number): string => {
  const now = Date.now();

  const diff = Math.floor((now - timestamp * 1000) / 1000);

  if (diff < 0) {
    return '0 second ago';
  }

  if (diff < 60) return `${diff} second${diff !== 1 ? 's' : ''} ago`;

  const minutes = Math.floor(diff / 60);
  if (minutes < 60) return `${minutes} minutes${minutes !== 1 ? 's' : ''} ago`;

  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours} hour${hours !== 1 ? 's' : ''} ago`;

  const days = Math.floor(hours / 24);
  if (days < 30) return `${days} day${days !== 1 ? 's' : ''} ago`;

  const date = new Date(timestamp * 1000);
  return date.toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  });
};

export const FormatSize = (size: number): string => {
  const kb = size / 1024;

  if (kb < 1024) {
    return `${FormatFloat(kb)} KB`;
  }

  const mb = kb / 1024;

  return `${FormatFloat(mb)} MB`;
};

export const FormatFloat = (value: number, decimals: number = 2) => {
  const factor = Math.pow(10, decimals);
  const truncated = Math.floor(value * factor) / factor;

  return truncated;
};

export const IsNumber = (n: string | null) => {
  try {
    if (!n) {
      return false;
    }

    return !!parseInt(n);
  } catch {
    return false;
  }
};
