import React, { Component } from 'react';
import PropTypes from 'prop-types';
import './InvoicePanel.css';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import BasicPanel from 'shared/BasicPanel';
import {
  makeGetUnbilledShipmentLineItems,
  makeTotalFromUnbilledLineItems,
} from 'shared/Entities/modules/shipmentLineItems';
import InvoiceTable from 'shared/Invoice/InvoiceTable';

export class InvoicePanel extends Component {
  constructor() {
    super();
    this.state = {};
  }

  render() {
    let invoicingContent = <span className="empty-content">No line items</span>;
    if (this.props.unbilledShipmentLineItems.length > 0) {
      invoicingContent = (
        <InvoiceTable
          shipmentLineItems={this.props.unbilledShipmentLineItems}
          totalAmount={this.props.lineItemsTotal}
        />
      );
    }

    return (
      <div className="invoice-panel">
        <BasicPanel title={'Invoicing'}>{invoicingContent}</BasicPanel>
      </div>
    );
  }
}

InvoicePanel.propTypes = {
  unbilledShipmentLineItems: PropTypes.array,
  shipmentId: PropTypes.string,
  lineItemsTotal: PropTypes.number,
};

//https://github.com/reduxjs/reselect#sharing-selectors-with-props-across-multiple-component-instances
const makeMapStateToProps = () => {
  //using a memoized selector
  const mapStateToProps = (state, ownProps) => {
    const getLineItems = makeGetUnbilledShipmentLineItems();
    const getLineItemSum = makeTotalFromUnbilledLineItems();
    return {
      unbilledShipmentLineItems: getLineItems(state, ownProps.shipmentId),
      lineItemsTotal: getLineItemSum(state, ownProps.shipmentId),
    };
  };
  return mapStateToProps;
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}
export default connect(makeMapStateToProps, mapDispatchToProps)(InvoicePanel);
