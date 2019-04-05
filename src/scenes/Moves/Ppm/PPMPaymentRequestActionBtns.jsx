import React from 'react';
import { withRouter } from 'react-router-dom';
import './PPMPaymentRequest.css';

const PPMPaymentRequestActionBtns = props => {
  const { nextBtnLabel, onClick, history } = props;
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
      <button onClick={onClick}>{nextBtnLabel}</button>
    </div>
  );
};

export default withRouter(PPMPaymentRequestActionBtns);
