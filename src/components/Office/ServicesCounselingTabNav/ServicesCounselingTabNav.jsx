import React from 'react';
import { NavLink } from 'react-router-dom';
import { Tag } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import 'styles/office.scss';
import TabNav from 'components/TabNav';

const ServicesCounselingTabNav = ({ unapprovedShipmentCount = 0, moveCode }) => {
  return (
    <header className="nav-header">
      <div className="grid-container-desktop-lg">
        <TabNav
          items={[
            <NavLink
              end
              className={({ isActive }) => (isActive ? 'usa-current' : '')}
              to={`/counseling/moves/${moveCode}/details`}
              data-testid="MoveDetails-Tab"
            >
              <span className="tab-title">Move details</span>
              {unapprovedShipmentCount > 0 && <Tag>{unapprovedShipmentCount}</Tag>}
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
              <span className="tab-title">Customer support remarks</span>
            </NavLink>,
            <NavLink
              end
              className={({ isActive }) => (isActive ? 'usa-current' : '')}
              to={`/counseling/moves/${moveCode}/history`}
              data-testid="MoveHistory-Tab"
            >
              <span className="tab-title">Move history</span>
            </NavLink>,
          ]}
        />
      </div>
    </header>
  );
};

ServicesCounselingTabNav.defaultProps = {
  unapprovedShipmentCount: 0,
};

ServicesCounselingTabNav.propTypes = {
  unapprovedShipmentCount: PropTypes.number,
  moveCode: PropTypes.string.isRequired,
};

export default ServicesCounselingTabNav;
