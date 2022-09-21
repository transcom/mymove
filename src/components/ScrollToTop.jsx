// https://github.com/ReactTraining/react-router/blob/master/packages/react-router-dom/docs/guides/scroll-restoration.md
// Use this component on any pages that should be auto-scrolled to the top on mount
// Should only be used once per-page and can be added at the App or global level if appropriate

import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';
// TODO: Make this function only take otherDeps, but change uses everywhere else
export default function ScrollToTop({ otherDep = null, otherDeps = [] }) {
  const { pathname } = useLocation();

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [pathname, otherDep, ...otherDeps]);

  return null;
}
