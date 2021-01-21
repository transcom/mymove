import React from 'react';
import { node, func } from 'prop-types';
import { Button, Title } from '@trussworks/react-uswds';

import MmLogo from '../../shared/images/milmove-logo.svg';

import styles from './index.module.scss';

import { CustomerShape } from 'types/moveOrder';

const MilMoveHeader = ({ children, customer, handleLogout }) => (
  <div className={styles.mmHeader}>
    <Title>
      <a href="/" title="office.move.mil" aria-label="office.move.mil">
        <img src={MmLogo} alt="MilMove Logo" />
      </a>
    </Title>
    <div className={styles.links}>
      {children}
      <div className={styles.verticalLine} />
      <span>
        {customer.last_name}, {customer.first_name}
      </span>
      <Button unstyled className={styles.signOut} onClick={handleLogout} type="button">
        Sign out
      </Button>
    </div>
  </div>
);

MilMoveHeader.propTypes = {
  children: node.isRequired,
  customer: CustomerShape.isRequired,
  handleLogout: func.isRequired,
};

export default MilMoveHeader;
