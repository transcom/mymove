import React from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import './PPMPaymentRequest.css';

const PPMPaymentRequestActionBtns = props => {
  const { nextBtnLabel, onClick, history, disabled, displaySaveForLater } = props;
  return (
    <div className="ppm-payment-request-footer">
      <div className="usa-width-two-thirds">
        <button
          type="button"
          className="usa-button-secondary"
          onClick={() => {
            history.push('/');
          }}
        >
          Cancel
        </button>
        {displaySaveForLater && (
          <button
            type="button"
            className="usa-button-secondary"
            onClick={() => {
              let result = Promise.resolve(onClick());
              result.then(value => {
                if (value === undefined) {
                  history.push('/');
                }
              });
            }}
            disabled={disabled}
          >
            Save For Later
          </button>
        )}
      </div>
      <button
        type="button"
        onClick={() => {
          onClick();
        }}
        disabled={disabled}
      >
        {nextBtnLabel}
      </button>
    </div>
  );
};
function mapStateToProps(state) {
  const { form } = state;
  let isDisabled = false;
  if (form.weight_ticket_wizard) {
    isDisabled = !(
      form.weight_ticket_wizard.values &&
      form.weight_ticket_wizard.values.vehicle_nickname &&
      form.weight_ticket_wizard.values.vehicle_options &&
      form.weight_ticket_wizard.values.empty_weight &&
      form.weight_ticket_wizard.values.full_weight &&
      form.weight_ticket_wizard.values.weight_ticket_date
    );
  }
  return {
    disabled: isDisabled,
  };
}

export default connect(mapStateToProps)(withRouter(PPMPaymentRequestActionBtns));
