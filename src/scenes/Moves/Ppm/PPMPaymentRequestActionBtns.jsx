import React from 'react';
import { withRouter } from 'react-router-dom';
import './PPMPaymentRequest.css';

const PPMPaymentRequestActionBtns = props => {
  const {
    nextBtnLabel,
    skipHandler,
    saveAndAddHandler,
    saveForLaterHandler,
    submitButtonsAreDisabled,
    displaySaveForLater,
    displaySkip,
    submitting,
    history,
  } = props;
  return (
    <div className="ppm-payment-request-footer">
      <div className="usa-width-two-thirds">
        <button type="button" className="usa-button-secondary" onClick={() => history.push('/')}>
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
      <div className="usa-width-one-third">
        {displaySkip && (
          <button data-cy="skip" type="button" className="usa-button-secondary" onClick={skipHandler}>
            Skip
          </button>
        )}
        <button type="button" onClick={saveAndAddHandler} disabled={submitButtonsAreDisabled || submitting}>
          {nextBtnLabel}
        </button>
      </div>
    </div>
  );
};

export default withRouter(PPMPaymentRequestActionBtns);
