import React from 'react';

import PPMPaymentRequestActionBtns from './PPMPaymentRequestActionBtns';
import './PPMPaymentRequest.css';

const PPMPaymentRequestIntro = () => {
  return (
    <div className="usa-grid">
      {/* TODO: change onclick handler to go to next page in flow */}
      <PPMPaymentRequestActionBtns onClick={() => {}} nextBtnLabel="Get Started" />
    </div>
  );
};
export default PPMPaymentRequestIntro;
