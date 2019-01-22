import React from 'react';
import PropTypes from 'prop-types';

import { withRouter } from 'react-router-dom';

const SaveCancelButtons = props => {
  const { submitting, valid } = props;
  const goBack = props.history.goBack;
  return (
    <div>
      <button type="submit" disabled={submitting || !valid}>
        Save
      </button>
      <button type="button" className="usa-button-secondary" disabled={submitting} onClick={goBack}>
        Cancel
      </button>
    </div>
  );
};

SaveCancelButtons.propTypes = {
  submitting: PropTypes.bool,
  valid: PropTypes.bool,
  history: PropTypes.shape({ goBack: PropTypes.func.isRequired }),
};

export default withRouter(SaveCancelButtons);
