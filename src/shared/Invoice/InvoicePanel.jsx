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

export class InvoicePanel extends Component {
  constructor() {
    super();
    this.state = {};
  }

  render() {
    let invoicingContent = <span className="empty-content">No line items</span>;
    if (this.props.shipmentLineItems.length > 0) {
      //stand up a table
    }

    return (
      <div className="invoice-panel">
        <BasicPanel title={'Invoicing'}>{invoicingContent}</BasicPanel>
      </div>
    );
  }
}

InvoicePanel.propTypes = {
  shipmentLineItems: PropTypes.array,
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
