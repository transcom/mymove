import React, { Component } from 'react';
import PropTypes from 'prop-types';
import './InvoicePanel.css';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import BasicPanel from 'shared/BasicPanel';
import { selectUnbilledShipmentLineItemsForShipment } from 'shared/Entities/modules/shipmentLineItems';
import { selectTariff400ngItems } from 'shared/Entities/modules/tariff400ngItems';

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
};

function mapStateToProps(state, ownProps) {
  return {
    shipmentLineItems: selectUnbilledShipmentLineItemsForShipment(state, ownProps.shipmentId),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators({}, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(InvoicePanel);
