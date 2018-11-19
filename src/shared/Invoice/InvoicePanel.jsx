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
import faClock from '@fortawesome/fontawesome-free-solid/faClock';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';

export class InvoicePanel extends Component {
  constructor() {
    super();
    this.state = {};
  }

  render() {
    let invoicingContent = <span className="empty-content">No line items</span>;
    let title = 'Invoicing';

    if (this.props.unbilledShipmentLineItems.length > 0) {
      invoicingContent = (
        <div>
          <InvoiceTable
            shipmentLineItems={this.props.unbilledShipmentLineItems}
            totalAmount={this.props.lineItemsTotal}
          />
        </div>
      );
      title = (
        <span>
          Invoicing{' '}
          <FontAwesomeIcon className="icon invoice-panel-icon-gold invoice-panel-icon--title" icon={faClock} />
        </span>
      );
    }

    return (
      <div className="invoice-panel">
        <BasicPanel title={title}>{invoicingContent}</BasicPanel>
      </div>
    );
  }
}

InvoicePanel.propTypes = {
  unbilledShipmentLineItems: PropTypes.array,
  shipmentId: PropTypes.string,
  shipmentStatus: PropTypes.string,
  lineItemsTotal: PropTypes.number,
};

//https://github.com/reduxjs/reselect#sharing-selectors-with-props-across-multiple-component-instances
const makeMapStateToProps = () => {
  //using a memoized selector
  const mapStateToProps = (state, ownProps) => {
    const getLineItems = makeGetUnbilledShipmentLineItems();
    const getLineItemSum = makeTotalFromUnbilledLineItems();
    const isShipmentDelivered = ownProps.shipmentStatus.toUpperCase() === 'DELIVERED';
    return {
      unbilledShipmentLineItems: isShipmentDelivered ? getLineItems(state, ownProps.shipmentId) : [],
      lineItemsTotal: isShipmentDelivered ? getLineItemSum(state, ownProps.shipmentId) : 0,
    };
  };
  return mapStateToProps;
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}
export default connect(makeMapStateToProps, mapDispatchToProps)(InvoicePanel);
