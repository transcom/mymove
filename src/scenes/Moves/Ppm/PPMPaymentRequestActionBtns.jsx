import React from 'react';
import { withRouter } from 'react-router-dom';
import './PPMPaymentRequest.css';
import Alert from 'shared/Alert';

const ConfirmationAlert = props => {
  const { confirmationAlertMsg, history } = props;
  return (
    <Alert type="warning" heading="">
      <div className="usa-width-two-thirds">{confirmationAlertMsg}</div>
      <div className="usa-width-one-thirds">
        <button type="button" className="usa-button-secondary" onClick={() => console.log('cancelled')}>
          Cancel
        </button>
        <button type="button" className="usa-button" onClick={() => history.push('/')}>
          OK
        </button>
      </div>
    </Alert>
  );
};

const PPMPaymentRequestActionBtns = props => {
  const {
    nextBtnLabel,
    displaySkip,
    skipHandler,
    saveAndAddHandler,
    finishLaterHandler,
    displayConfirmation,
    submitButtonsAreDisabled,
    submitting,
  } = props;
  return (
    <div className="ppm-payment-request-footer">
      <div className="usa-width-one-whole">
        <ConfirmationAlert confirmationAlertMsg="Partially completed entries will not be saved. Click OK to continue. Click cancel to return and edit." />
      </div>
      {displayConfirmation && (
        <div>
          <div className="usa-width-two-thirds">
            <button
              type="button"
              className="usa-button-secondary"
              onClick={finishLaterHandler}
              disabled={submitButtonsAreDisabled || submitting}
            >
              Finish Later
            </button>
          </div>
          <div className="usa-width-one-thirds">
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
      )}
    </div>
  );
};

export default withRouter(PPMPaymentRequestActionBtns);
