import React from 'react';
import { node, func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { NavLink } from 'react-router-dom';

import MmLogo from '../../shared/images/milmove-logo.svg';

import styles from './index.module.scss';

import { CustomerShape } from 'types/moveOrder';

const MilMoveHeader = ({ children, customer, handleLogout }) => (
  <div className={styles.mmHeader}>
    <div className="usa-logo" id="basic-logo">
      <em className="usa-logo__text">
        <NavLink to="/" title="office.move.mil" aria-label="office.move.mil">
          <img src={MmLogo} alt="MilMove Logo" />
        </NavLink>
      </em>
    </div>
    <div className={styles.links}>
      {children}
      <span className={styles.lineAdd}>&nbsp;</span>
      <span>
        {customer.last_name}, {customer.first_name}
      </span>
      <span>
        <Button className={styles.signOut} disabled={false} onClick={handleLogout}>
          <span>Sign out</span>
        </Button>
      </span>
    </div>
  </div>
);

MilMoveHeader.propTypes = {
  children: node.isRequired,
  customer: CustomerShape.isRequired,
  handleLogout: func.isRequired,
};

export default MilMoveHeader;
