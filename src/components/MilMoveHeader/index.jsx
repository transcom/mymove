import React from 'react';
import { node, func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import { ReactComponent as MmLogo } from '../../shared/images/milmove-logo.svg';

import styles from './index.module.scss';

import { CustomerShape } from 'types/moveOrder';

const MilMoveHeader = ({ children, customer, handleLogout }) => (
  <div className={styles.mmHeader}>
    <MmLogo />
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
