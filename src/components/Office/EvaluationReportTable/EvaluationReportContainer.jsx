import React from 'react';
import * as PropTypes from 'prop-types';

const EvaluationReportContainer = ({ className, children }) => {
  return (
    <div className={className} data-testid="EvaluationReportContainer">
      {children}
    </div>
  );
};

EvaluationReportContainer.propTypes = {
  className: PropTypes.string,
  children: PropTypes.node.isRequired,
};

EvaluationReportContainer.defaultProps = {
  className: '',
};

export default EvaluationReportContainer;
