import React from 'react';
import PropTypes from 'prop-types';
import withRouter from 'utils/routing';

const SaveCancelButtons = (props) => {
  const {
    submitting,
    valid,
    router: { navigate },
  } = props;
  return (
    <div className="margin-top-2">
      <button className="usa-button margin-bottom-1" type="submit" disabled={submitting || !valid}>
        Save
      </button>
      <button
        type="button"
        className="usa-button usa-button--secondary margin-bottom-1"
        disabled={submitting}
        onClick={() => navigate(-1)}
      >
        Cancel
      </button>
    </div>
  );
};

SaveCancelButtons.propTypes = {
  submitting: PropTypes.bool,
  valid: PropTypes.bool,
  router: PropTypes.object,
};

export default withRouter(SaveCancelButtons);
