import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import Alert from 'shared/Alert';
import './InvoicePanel.css';

class InvoicePaymentAlert extends PureComponent {
  render() {
    let paymentAlert;
    const status = this.props.createInvoiceStatus;

    if (status.error) {
      paymentAlert = (
        <Alert type="error" heading="Oops, something went wrong!">
          <span className="warning--header">Please try again.</span>
        </Alert>
      );
    } else if (status.isLoading) {
      paymentAlert = (
        <Alert type="loading" heading="Creating invoice">
          <span className="warning--header">Sending information to USBank/Syncada.</span>
        </Alert>
      );
    } else if (status.isSuccess) {
      paymentAlert = (
        <div>
          <Alert type="success" heading="Success!">
            <span className="warning--header">The invoice has been created and will be paid soon.</span>
          </Alert>
        </div>
      );
    }

    return <div>{paymentAlert}</div>;
  }
}

InvoicePaymentAlert.propTypes = {
  createInvoiceStatus: PropTypes.object,
};

export default InvoicePaymentAlert;
