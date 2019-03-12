import React, { Fragment } from 'react';
import Alert from 'shared/Alert';

export const Code35FormAlert = props => {
  return (
    <Fragment>
      {props.showAlert && (
        <Alert type="warning" heading="Oops– you appear to have exceeded your estimated cost!">
          <span className="warning--header">
            Amount exceeds approved estimate. If you continue, you'll only be paid the max approved amount. Submit a
            separate pre-approval request to cover any additional costs.
          </span>
        </Alert>
      )}
    </Fragment>
  );
};
