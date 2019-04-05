import React from 'react';

import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import './PPMPaymentRequest.css';

const PPMPaymentRequestIntro = props => {
  const { history, match } = props;
  return (
    <div className="usa-grid">
      <PPMPaymentRequestActionBtns
        onClick={() => {
          history.push(`/moves/${match.params.moveId}/ppm-weight-ticket`);
        }}
        nextBtnLabel="Get Started"
      />
    </div>
  );
};
export default PPMPaymentRequestIntro;
