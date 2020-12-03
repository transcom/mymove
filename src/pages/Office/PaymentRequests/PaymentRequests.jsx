import React from 'react';
import { withRouter } from 'react-router-dom';
import { Button } from '@trussworks/react-uswds';

import { MatchShape, HistoryShape } from 'types/router';

const PaymentRequests = ({ history, match }) => {
  const { moveOrderId } = match.params;
  const paymentRequestId = history.location.state.detail;

  const handleClick = () => {
    history.push(`/moves/${moveOrderId}/payment-requests/${paymentRequestId}`);
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

PaymentRequests.propTypes = {
  history: HistoryShape.isRequired,
  match: MatchShape.isRequired,
};

export default withRouter(PaymentRequests);
