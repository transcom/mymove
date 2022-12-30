import React from 'react';
import { useNavigate } from 'react-router-dom';

function CustomerAgreementLegalese(props) {
  const navigate = useNavigate();
  function goBack() {
    navigate(-1);
  }

  return (
    <div className="grid-container usa-prose customer-agreement-legalese-container">
      <div className="grid-row">
        <div className="grid-col-12">
          <div>
            <a onClick={goBack} className="usa-link">
              {'<'} Back
            </a>
          </div>
          <h1>Customer Agreement</h1>
          <p style={{ marginBottom: '20px' }}>
            Before submitting your payment request, please carefully read the following:
          </p>
          <div className="usa-grid-full customer-agreement-legalese-text">
            <h4>LEGAL AGREEMENT / PRIVACY ACT</h4>
            <h4>Financial Liability:</h4>
            If this shipment(s) incurs costs above the allowance I am entitled to, I will pay the difference to the
            government, or consent to the collection from my pay as necessary to cover all excess costs associated by
            this shipment(s).
            <h4>Advance Obligation:</h4>
            <p>
              I understand that the maximum advance allowed is based on the estimated weight and scheduled departure
              date of my shipment(s). In the event, less weight is moved or my move occurs on a different scheduled
              departure date, I may have to remit the difference with the balance of my incentive disbursement and/or
              from the collection of my pay as may be necessary.
            </p>
            <p>
              I understand that the maximum advance allowed is based on the estimated weight and scheduled departure
              date of my shipment(s). In the event, less weight is moved or my move occurs on a different scheduled
              departure date, I may have to remit the difference with the balance of my incentive disbursement and/or
              from the collection of my pay as may be necessary. If I receive an advance for my PPM shipment, I agree to
              furnish weight tickets within 45 days of final delivery to my destination. I understand that failure to
              furnish weight tickets within this timeframe may lead to the collection of my pay as necessary to cover
              the cost of the advance.
            </p>
          </div>
          <div className="usa-grid button-bar">
            <button className="usa-button usa-button--secondary" data-testid="back-button" onClick={goBack}>
              Back
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export default CustomerAgreementLegalese;
