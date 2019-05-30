import React from 'react';
import { withRouter } from 'react-router-dom';
import './PPMPaymentRequest.css';

const PPMPaymentRequestActionBtns = props => {
  const {
    nextBtnLabel,
    cancelHandler,
    saveAndAddHandler,
    saveForLaterHandler,
    submitButtonsAreDisabled,
    displaySaveForLater,
    submitting,
  } = props;
  return (
    <div className="ppm-payment-request-footer">
      <div className="usa-width-two-thirds">
        <button type="button" className="usa-button-secondary" onClick={cancelHandler}>
          Cancel
        </button>
        {displaySaveForLater && (
          <button
            type="button"
            className="usa-button-secondary"
            onClick={saveForLaterHandler}
            disabled={submitButtonsAreDisabled || submitting}
          >
            Save For Later
          </button>
        )}
      </div>
      <button type="button" onClick={saveAndAddHandler} disabled={submitButtonsAreDisabled || submitting}>
        {nextBtnLabel}
      </button>
    </div>
  );
};

export default withRouter(PPMPaymentRequestActionBtns);
