import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';

import { isOfficeSite } from 'shared/constants.js';
import LineItemTable from 'shared/Invoice/LineItemTable';
import Alert from 'shared/Alert';
import { isLoading } from 'shared/constants';

import styles from './UnbilledTable.module.scss';
import panelStyles from './InvoicePanel.module.scss';

export class UnbilledTable extends PureComponent {
  constructor(props) {
    super(props);

    this.state = {
      draftInvoice: false,
    };
  }

  draftInvoice = () => {
    this.setState({ draftInvoice: true });
  };

  cancelPayment = () => {
    this.setState({ draftInvoice: false });
  };

  approvePayment = () => {
    this.props.approvePayment();
    this.cancelPayment();
  };

  render() {
    const allowPayments =
      this.props.allowPayments &&
      isOfficeSite && //user is an office user
      this.props.createInvoiceStatus !== isLoading;

    // Table header, also contains buttons for initiating invoice payment
    let header;
    if (this.state.draftInvoice) {
      header = (
        <Alert type="warning" heading="Approve payment?">
          <span className={panelStyles['warning--header']}>Please make sure you've double-checked everything.</span>
          <button className="button usa-button-secondary" onClick={this.cancelPayment}>
            Cancel
          </button>
          <button className="button usa-button-primary" onClick={this.approvePayment}>
            Approve
          </button>
        </Alert>
      );
    } else {
      header = (
        <div className={panelStyles['invoice-panel__header-cont']} data-cy="invoice-panel__header-cont">
          <div className="usa-width-one-half">
            <h5>Unbilled line items</h5>
          </div>
          {allowPayments && (
            <div className={`usa-width-one-half ${styles['button__approve-payment']}`}>
              <button className="button button-secondary" onClick={this.draftInvoice}>
                Approve Payment
              </button>
            </div>
          )}
        </div>
      );
    }

    let itemsComponent;
    if (this.props.lineItems.length) {
      itemsComponent = (
        <LineItemTable shipmentLineItems={this.props.lineItems} totalAmount={this.props.lineItemsTotal} />
      );
    }

    return (
      <div className={panelStyles['invoice-panel__table-cont']} data-cy="unbilled-table">
        {header}
        {itemsComponent}
      </div>
    );
  }
}

UnbilledTable.propTypes = {
  lineItems: PropTypes.array,
  lineItemsTotal: PropTypes.number,
  approvePayment: PropTypes.func,
  allowPayments: PropTypes.bool,
  createInvoiceStatus: PropTypes.string,
};

export default UnbilledTable;
