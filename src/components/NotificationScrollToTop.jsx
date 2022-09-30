import { useEffect, useRef } from 'react';

export default function NotificationScrollToTop({ dependency }) {
  const isMounted = useRef(false);

  useEffect(() => {
    if (isMounted.current) {
      window.scrollTo(0, 0);
    } else {
      isMounted.current = true;
    }
  }, [dependency]);

  return null;
}
