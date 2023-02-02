import { useEffect, useRef } from 'react';
import { object } from 'prop-types';

export default function NotificationScrollToTop({ dependency, target }) {
  const isMounted = useRef(false);

  useEffect(() => {
    if (isMounted.current) {
      if (target) {
        target.scrollTo(0, 0);
      }
    } else {
      isMounted.current = true;
    }
  }, [dependency, target]);

  return null;
}

NotificationScrollToTop.propTypes = {
  target: object,
};

NotificationScrollToTop.defaultProps = {
  target: window,
};
