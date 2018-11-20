import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';
import './InvoicePanel.css';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import { isOfficeSite, isDevelopment } from 'shared/constants.js';

import BasicPanel from 'shared/BasicPanel';
import {
  makeGetUnbilledShipmentLineItems,
  makeTotalFromUnbilledLineItems,
} from 'shared/Entities/modules/shipmentLineItems';

export class InvoicePanel extends PureComponent {
  approvePayment = () => {
    //this.props.onApprovePayment(this.props.shipmentId);
    console.log('Approve Payment button clicked!');
  };

  render() {
    let invoicingContent = <div className="empty-content">No line items</div>;
    if (this.props.shipmentLineItems.length > 0) {
      //stand up a table
    }

    return (
      <div className="invoice-panel">
        <BasicPanel title={'Invoicing'}>
          {isOfficeSite &&
            this.props.shipmentState === 'DELIVERED' && (
              <div className="usa-width-one-whole align-right">
                <button
                  className="button button-secondary"
                  disabled={!this.props.canApprove && !isDevelopment}
                  onClick={this.approvePayment}
                >
                  Approve Payment
                </button>
              </div>
            )}
          {invoicingContent}
        </BasicPanel>
      </div>
    );
  }
}

InvoicePanel.propTypes = {
  shipmentLineItems: PropTypes.array,
  shipmentId: PropTypes.string,
  shipmentState: PropTypes.string,
  lineItemsTotal: PropTypes.number,
  onApprovePayment: PropTypes.func,
  canApprove: PropTypes.bool,
};

//https://github.com/reduxjs/reselect#sharing-selectors-with-props-across-multiple-component-instances
const makeMapStateToProps = () => {
  //using a memoized selector
  const mapStateToProps = (state, ownProps) => {
    const getLineItems = makeGetUnbilledShipmentLineItems();
    const getLineItemSum = makeTotalFromUnbilledLineItems();
    return {
      shipmentLineItems: getLineItems(state, ownProps.shipmentId),
      lineItemsTotal: getLineItemSum(state, ownProps.shipmentId),
    };
  };
  return mapStateToProps;
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}
export default connect(makeMapStateToProps, mapDispatchToProps)(InvoicePanel);
