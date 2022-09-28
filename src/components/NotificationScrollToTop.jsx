import { useEffect } from 'react';

export default function NotificationScrollToTop({ dependency }) {
  useEffect(() => {
    window.scrollTo(0, 0);
  }, [dependency]);

  return null;
}
