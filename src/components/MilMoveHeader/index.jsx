import React from 'react';
import { func, string, node } from 'prop-types';
import { Button, Title } from '@trussworks/react-uswds';

import MmLogo from '../../shared/images/milmove-logo.svg';

import styles from './index.module.scss';

const MilMoveHeader = ({ children, handleLogout, firstName, lastName }) => (
  <div className={styles.mmHeader}>
    <Title>
      <a href="/" title="office.move.mil" aria-label="office.move.mil">
        <img src={MmLogo} alt="MilMove Logo" />
      </a>
    </Title>
    <div className={styles.links}>
      {children}
      <div className={styles.verticalLine} />
      {lastName !== '' && firstName !== '' && (
        <span>
          {lastName}, {firstName}
        </span>
      )}
      <Button unstyled className={styles.signOut} onClick={handleLogout} type="button">
        Sign out
      </Button>
    </div>
  </div>
);

MilMoveHeader.defaultProps = {
  firstName: '',
  lastName: '',
};

MilMoveHeader.propTypes = {
  children: node.isRequired,
  firstName: string,
  lastName: string,
  handleLogout: func.isRequired,
};

export default MilMoveHeader;
