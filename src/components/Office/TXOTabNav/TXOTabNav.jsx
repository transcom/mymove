import React from 'react';
import { NavLink } from 'react-router-dom';
import { Tag } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import 'styles/office.scss';
import TabNav from 'components/TabNav';
import { OrdersShape } from 'types/customerShapes';

const TXOTabNav = ({
  unapprovedShipmentCount,
  unapprovedServiceItemCount,
  excessWeightRiskCount,
  pendingPaymentRequestCount,
  unapprovedSITExtensionCount,
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

  let moveTaskOrderTagCount = 0;
  if (unapprovedServiceItemCount > 0) {
    moveTaskOrderTagCount += unapprovedServiceItemCount;
  }
  if (excessWeightRiskCount > 0) {
    moveTaskOrderTagCount += 1;
  }
  if (unapprovedSITExtensionCount > 0) {
    moveTaskOrderTagCount += unapprovedSITExtensionCount;
  }

  return (
    <header className="nav-header">
      <div className="grid-container-desktop-lg">
        <TabNav
          items={[
            <NavLink
              exact
              activeClassName="usa-current"
              to={`/moves/${moveCode}/details`}
              data-testid="MoveDetails-Tab"
            >
              <span className="tab-title">Move details</span>
              {moveDetailsTagCount > 0 && <Tag>{moveDetailsTagCount}</Tag>}
            </NavLink>,
            <NavLink data-testid="MoveTaskOrder-Tab" exact activeClassName="usa-current" to={`/moves/${moveCode}/mto`}>
              <span className="tab-title">Move task order</span>
              {moveTaskOrderTagCount > 0 && <Tag>{moveTaskOrderTagCount}</Tag>}
            </NavLink>,
            <NavLink exact activeClassName="usa-current" to={`/moves/${moveCode}/payment-requests`}>
              <span className="tab-title">Payment requests</span>
              {pendingPaymentRequestCount > 0 && <Tag>{pendingPaymentRequestCount}</Tag>}
            </NavLink>,
            <NavLink exact activeClassName="usa-current" to={`/moves/${moveCode}/customer-support-remarks`}>
              <span className="tab-title">Customer support remarks</span>
            </NavLink>,
            <NavLink exact activeClassName="usa-current" to={`/moves/${moveCode}/evaluation-reports`}>
              <span className="tab-title">Evaluation reports</span>
            </NavLink>,
            <NavLink exact activeClassName="usa-current" to={`/moves/${moveCode}/history`}>
              <span className="tab-title">Move history</span>
            </NavLink>,
          ]}
        />
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
};

TXOTabNav.propTypes = {
  order: OrdersShape.isRequired,
  unapprovedShipmentCount: PropTypes.number,
  unapprovedServiceItemCount: PropTypes.number,
  excessWeightRiskCount: PropTypes.number,
  pendingPaymentRequestCount: PropTypes.number,
  unapprovedSITExtensionCount: PropTypes.number,
  moveCode: PropTypes.string.isRequired,
};

export default TXOTabNav;
