import React from 'react';
import { NavLink } from 'react-router-dom';
import { Tag } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './TXOTabNav.module.scss';

import 'styles/office.scss';
import TabNav from 'components/TabNav';
import { OrdersShape } from 'types/customerShapes';

const TXOTabNav = ({
  unapprovedShipmentCount,
  unapprovedServiceItemCount,
  excessWeightRiskCount,
  pendingPaymentRequestCount,
  unapprovedSITExtensionCount,
  missingOrdersInfoCount,
  shipmentErrorConcernCount,
  shipmentsWithDeliveryAddressUpdateRequestedCount,
  order,
  moveCode,
}) => {
  let moveDetailsTagCount = 0;
  if (unapprovedShipmentCount > 0) {
    moveDetailsTagCount += unapprovedShipmentCount;
  }
  if (order.uploadedAmendedOrderID && !order.amendedOrdersAcknowledgedAt) {
    moveDetailsTagCount += 1;
  }
  if (shipmentErrorConcernCount) {
    moveDetailsTagCount += shipmentErrorConcernCount;
  }
  if (shipmentsWithDeliveryAddressUpdateRequestedCount) {
    moveDetailsTagCount += shipmentsWithDeliveryAddressUpdateRequestedCount;
  }
  if (missingOrdersInfoCount > 0) {
    moveDetailsTagCount += missingOrdersInfoCount;
  }

  let moveTaskOrderTagCount = 0;
  if (unapprovedServiceItemCount > 0) {
    moveTaskOrderTagCount += unapprovedServiceItemCount;
  }
  if (excessWeightRiskCount > 0) {
    moveTaskOrderTagCount += excessWeightRiskCount;
  }
  if (unapprovedSITExtensionCount > 0) {
    moveTaskOrderTagCount += unapprovedSITExtensionCount;
  }

  const items = [
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/moves/${moveCode}/details`}
      data-testid="MoveDetails-Tab"
    >
      <span className="tab-title">Move Details</span>
      {moveDetailsTagCount > 0 && <Tag>{moveDetailsTagCount}</Tag>}
    </NavLink>,
    <NavLink
      data-testid="MoveTaskOrder-Tab"
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/moves/${moveCode}/mto`}
    >
      <span className="tab-title">Move Task Order</span>
      {moveTaskOrderTagCount > 0 && <Tag>{moveTaskOrderTagCount}</Tag>}
    </NavLink>,
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/moves/${moveCode}/payment-requests`}
    >
      <span className="tab-title">Payment Requests</span>
      {pendingPaymentRequestCount > 0 && <Tag>{pendingPaymentRequestCount}</Tag>}
    </NavLink>,
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/moves/${moveCode}/customer-support-remarks`}
    >
      <span className="tab-title">Customer Support Remarks</span>
    </NavLink>,
    <NavLink className={({ isActive }) => (isActive ? 'usa-current' : '')} to={`/moves/${moveCode}/evaluation-reports`}>
      <span className="tab-title">Quality Assurance</span>
    </NavLink>,
    <NavLink end className={({ isActive }) => (isActive ? 'usa-current' : '')} to={`/moves/${moveCode}/history`}>
      <span className="tab-title">Move History</span>
    </NavLink>,
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to="supporting-documents"
      data-testid="SupportingDocuments-Tab"
    >
      <span className="tab-title">Supporting Documents</span>
    </NavLink>,
  ];

  return (
    <header className="nav-header">
      <div className={classnames('grid-container-desktop-lg', styles.TabNav)}>
        <TabNav items={items} />
      </div>
    </header>
  );
};

TXOTabNav.defaultProps = {
  unapprovedShipmentCount: 0,
  unapprovedServiceItemCount: 0,
  excessWeightRiskCount: 0,
  pendingPaymentRequestCount: 0,
  unapprovedSITExtensionCount: 0,
  shipmentsWithDeliveryAddressUpdateRequestedCount: 0,
};

TXOTabNav.propTypes = {
  order: OrdersShape.isRequired,
  unapprovedShipmentCount: PropTypes.number,
  unapprovedServiceItemCount: PropTypes.number,
  excessWeightRiskCount: PropTypes.number,
  pendingPaymentRequestCount: PropTypes.number,
  unapprovedSITExtensionCount: PropTypes.number,
  shipmentsWithDeliveryAddressUpdateRequestedCount: PropTypes.number,
  moveCode: PropTypes.string.isRequired,
};

export default TXOTabNav;
