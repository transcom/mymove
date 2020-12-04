import React from 'react';
import { withRouter } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';

import { MatchShape, HistoryShape } from 'types/router';

const PaymentRequests = ({ history, match }) => {
  const { paymentRequestID, moveLocator } = match.params;
  // values.id = payment request id

  const handleClick = () => {
    history.push(`/moves/${moveLocator}/payment-requests/${paymentRequestID}`);
  };

  return (
    <div className="grid-container-desktop-lg" data-testid="PaymentRequests">
      <h1>Payment requests</h1>
      <div className="container">
        <Button data-testid="ReviewServiceItems" onClick={handleClick}>
          Review service items
        </Button>
      </div>
    </div>
  );
};

// include an array list to iterate over for all Payment Requests in the move
// tie appearance of button to status. status == 'APPROVED' && (render the card)

PaymentRequests.propTypes = {
  history: HistoryShape.isRequired,
  match: MatchShape.isRequired,
};

export default withRouter(PaymentRequests);
