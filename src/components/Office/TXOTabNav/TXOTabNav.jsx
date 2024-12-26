import React from 'react';
import { NavLink } from 'react-router-dom';
import { Tag } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './TXOTabNav.module.scss';

import 'styles/office.scss';
import TabNav from 'components/TabNav';
import { OrdersShape } from 'types/customerShapes';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

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
  const [supportingDocsFF, setSupportingDocsFF] = React.useState(false);
  React.useEffect(() => {
    const fetchData = async () => {
      setSupportingDocsFF(await isBooleanFlagEnabled('manage_supporting_docs'));
    };
    fetchData();
  }, []);

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
      <span className="tab-title">Move details</span>
      {moveDetailsTagCount > 0 && <Tag>{moveDetailsTagCount}</Tag>}
    </NavLink>,
    <NavLink
      data-testid="MoveTaskOrder-Tab"
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/moves/${moveCode}/mto`}
    >
      <span className="tab-title">Move task order</span>
      {moveTaskOrderTagCount > 0 && <Tag>{moveTaskOrderTagCount}</Tag>}
    </NavLink>,
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/moves/${moveCode}/payment-requests`}
    >
      <span className="tab-title">Payment requests</span>
      {pendingPaymentRequestCount > 0 && <Tag>{pendingPaymentRequestCount}</Tag>}
    </NavLink>,
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/moves/${moveCode}/customer-support-remarks`}
    >
      <span className="tab-title">Customer support remarks</span>
    </NavLink>,
    <NavLink className={({ isActive }) => (isActive ? 'usa-current' : '')} to={`/moves/${moveCode}/evaluation-reports`}>
      <span className="tab-title">Quality assurance</span>
    </NavLink>,
    <NavLink end className={({ isActive }) => (isActive ? 'usa-current' : '')} to={`/moves/${moveCode}/history`}>
      <span className="tab-title">Move history</span>
    </NavLink>,
  ];

  if (supportingDocsFF)
    items.push(
      <NavLink
        end
        className={({ isActive }) => (isActive ? 'usa-current' : '')}
        to="supporting-documents"
        data-testid="SupportingDocuments-Tab"
      >
        <span className="tab-title">Supporting Documents</span>
      </NavLink>,
    );

  return (
    <header className="nav-header">
      <div
        className={
          supportingDocsFF ? classnames('grid-container-desktop-lg', styles.TabNav) : 'grid-container-desktop-lg'
        }
      >
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
