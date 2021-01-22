import React from 'react';
import { isEmpty } from 'lodash';
import { func, node } from 'prop-types';
import { Button, Title } from '@trussworks/react-uswds';

import MmLogo from '../../shared/images/milmove-logo.svg';

import styles from './index.module.scss';

import { OfficeUserInfoShape } from 'types/index';

const MilMoveHeader = ({ children, handleLogout, officeUser }) => (
  <div className={styles.mmHeader}>
    <Title>
      <a href="/" title="office.move.mil" aria-label="office.move.mil">
        <img src={MmLogo} alt="MilMove Logo" />
      </a>
    </Title>
    <div className={styles.links}>
      {children}
      <div className={styles.verticalLine} />
      {!isEmpty(officeUser) && (
        <span>
          {officeUser.last_name}, {officeUser.first_name}
        </span>
      )}
      <Button unstyled className={styles.signOut} onClick={handleLogout} type="button">
        Sign out
      </Button>
    </div>
  </div>
);

MilMoveHeader.propTypes = {
  children: node.isRequired,
  officeUser: OfficeUserInfoShape.isRequired,
  handleLogout: func.isRequired,
};

export default MilMoveHeader;
