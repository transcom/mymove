import React from 'react';

function FeedbackConfirmation({ confirmationText, pendingValue }) {
  return (
    <div>
      <p>{confirmationText}</p>
      <p>{pendingValue}</p>
    </div>
  );
}

export default FeedbackConfirmation;
