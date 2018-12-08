import React, { PureComponent } from 'react';
import PropTypes from 'prop-types';
import './InvoicePanel.css';

import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';
import BasicPanel from 'shared/BasicPanel';
import {
  selectUnbilledShipmentLineItems,
  selectTotalFromUnbilledLineItems,
} from 'shared/Entities/modules/shipmentLineItems';
import { sendHHGInvoice } from 'scenes/Office/ducks';
import InvoiceTable from 'shared/Invoice/InvoiceTable';

export class InvoicePanel extends PureComponent {
  approvePayment = () => {
    return this.props.sendHHGInvoice(this.props.shipmentId);
  };

  render() {
    let invoicingContent = <span className="empty-content">No line items</span>;
    let title = 'Invoicing';

    if (this.props.unbilledShipmentLineItems.length > 0) {
      invoicingContent = (
        <div>
          <InvoiceTable
            shipmentLineItems={this.props.unbilledShipmentLineItems}
            shipmentStatus={this.props.shipmentStatus}
            totalAmount={this.props.lineItemsTotal}
            approvePayment={this.approvePayment}
            canApprove={this.props.canApprove}
          />
        </div>
      );
      title = <span>Invoicing</span>;
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
  onApprovePayment: PropTypes.func,
  canApprove: PropTypes.bool,
};

const mapStateToProps = (state, ownProps) => {
  const isShipmentDelivered = ownProps.shipmentStatus.toUpperCase() === 'DELIVERED';
  return {
    unbilledShipmentLineItems: isShipmentDelivered ? selectUnbilledShipmentLineItems(state, ownProps.shipmentId) : [],
    lineItemsTotal: isShipmentDelivered ? selectTotalFromUnbilledLineItems(state, ownProps.shipmentId) : 0,
  };
};

function mapDispatchToProps(dispatch) {
  return bindActionCreators({ sendHHGInvoice }, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(InvoicePanel);
