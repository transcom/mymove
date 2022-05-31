// https://github.com/ReactTraining/react-router/blob/master/packages/react-router-dom/docs/guides/scroll-restoration.md
// Use this component on any pages that should be auto-scrolled to the top on mount
// Should only be used once per-page and can be added at the App or global level if appropriate

import { useEffect } from 'react';
import { useLocation } from 'react-router-dom';
import PropTypes from 'prop-types';

const ScrollToTop = ({ otherDep = null }) => {
  const { pathname } = useLocation();

  useEffect(() => {
    window.scrollTo(0, 0);
  }, [pathname, otherDep]);

  return null;
};

ScrollToTop.propTypes = {
  otherDep: PropTypes.oneOfType([PropTypes.func, PropTypes.node]),
};

ScrollToTop.defaultProps = {
  otherDep: null,
};

export default ScrollToTop;
