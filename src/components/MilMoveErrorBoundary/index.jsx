import React from 'react';
import * as PropTypes from 'prop-types';

import { milmoveLogger } from 'utils/milmoveLog';
import { retryPageLoading } from 'utils/retryPageLoading';

// This error boundary will probably not be reached for most errors.
// See more comments in App/index.jsx
export class MilMoveErrorBoundary extends React.Component {
  constructor(props) {
    super(props);
    this.state = { hasError: false };
  }

  componentDidCatch(error, info) {
    this.setState({ hasError: true });
    const { message } = error;
    milmoveLogger.error({ message, info });
    retryPageLoading(error);
  }

  render() {
    const { hasError } = this.state;
    const { children, fallback } = this.props;
    if (hasError) {
      return fallback;
    }
    return children;
  }
}

MilMoveErrorBoundary.propTypes = {
  children: PropTypes.node.isRequired,
  fallback: PropTypes.node.isRequired,
};

export default MilMoveErrorBoundary;
