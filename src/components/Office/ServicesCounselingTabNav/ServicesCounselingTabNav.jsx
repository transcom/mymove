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
              exact
              activeClassName="usa-current"
              to={`/counseling/moves/${moveCode}/details`}
              data-testid="MoveDetails-Tab"
            >
              <span className="tab-title">Move details</span>
              {unapprovedShipmentCount > 0 && <Tag>{unapprovedShipmentCount}</Tag>}
            </NavLink>,
            <NavLink exact activeClassName="usa-current" to={`/counseling/moves/${moveCode}/customer-support-remarks`}>
              <span className="tab-title">Customer support remarks</span>
            </NavLink>,
            <NavLink
              exact
              activeClassName="usa-current"
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
