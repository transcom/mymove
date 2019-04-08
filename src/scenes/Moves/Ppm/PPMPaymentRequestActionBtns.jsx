import React from 'react';
import { withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import './PPMPaymentRequest.css';

const PPMPaymentRequestActionBtns = props => {
  const { nextBtnLabel, onClick, history, disabled } = props;
  return (
    <div className="ppm-payment-request-footer">
      <button
        className="usa-button-secondary"
        onClick={() => {
          history.push('/');
        }}
      >
        Cancel
      </button>
      <button onClick={onClick} disabled={disabled}>
        {nextBtnLabel}
      </button>
    </div>
  );
};
function mapStateToProps(state) {
  const { form } = state;
  let isDisabled = true;
  if (form.weight_ticket_wizard && form.weight_ticket_wizard.values) {
    isDisabled = !(
      form.weight_ticket_wizard.values.vehicle_nickname && form.weight_ticket_wizard.values.vehicle_options
    );
  }
  return {
    disabled: isDisabled,
  };
}

export default connect(mapStateToProps)(withRouter(PPMPaymentRequestActionBtns));
