import React from 'react';
import { NavLink } from 'react-router-dom';
import { Tag } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './ServicesCounselingTabNav.module.scss';

import 'styles/office.scss';
import TabNav from 'components/TabNav';

const ServicesCounselingTabNav = ({
  shipmentWarnConcernCount = 0,
  shipmentErrorConcernCount,
  missingOrdersInfoCount,
  moveCode,
}) => {
  let moveDetailsTagCount = 0;
  if (shipmentWarnConcernCount > 0) {
    moveDetailsTagCount += shipmentWarnConcernCount;
  }
  if (shipmentErrorConcernCount > 0) {
    moveDetailsTagCount += shipmentErrorConcernCount;
  }
  if (missingOrdersInfoCount > 0) {
    moveDetailsTagCount += missingOrdersInfoCount;
  }

  const items = [
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/counseling/moves/${moveCode}/details`}
      data-testid="MoveDetails-Tab"
    >
      <span className="tab-title">Move Details</span>
      {moveDetailsTagCount > 0 && <Tag>{moveDetailsTagCount}</Tag>}
    </NavLink>,
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/counseling/moves/${moveCode}/mto`}
      data-testid="MoveTaskOrder-Tab"
    >
      <span className="tab-title">Move Task Order</span>
    </NavLink>,
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/counseling/moves/${moveCode}/customer-support-remarks`}
    >
      <span className="tab-title">Customer Support Remarks</span>
    </NavLink>,
    <NavLink
      end
      className={({ isActive }) => (isActive ? 'usa-current' : '')}
      to={`/counseling/moves/${moveCode}/history`}
      data-testid="MoveHistory-Tab"
    >
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
        <div id="shipments" />
      </div>
    </header>
  );
};

ServicesCounselingTabNav.defaultProps = {};

ServicesCounselingTabNav.propTypes = {
  moveCode: PropTypes.string.isRequired,
};

export default ServicesCounselingTabNav;
