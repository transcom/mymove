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
import InvoiceTable from 'shared/Invoice/InvoiceTable';
import faClock from '@fortawesome/fontawesome-free-solid/faClock';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';

export class InvoicePanel extends PureComponent {
  approvePayment = () => {
    //this.props.onApprovePayment(this.props.shipmentId);
    console.log('Approve Payment button clicked!');
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
  return bindActionCreators({}, dispatch);
}
export default connect(mapStateToProps, mapDispatchToProps)(InvoicePanel);
