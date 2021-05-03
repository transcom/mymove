import React from 'react';
import { node } from 'prop-types';
import { Header, Title } from '@trussworks/react-uswds';

import MmLogo from '../../shared/images/milmove-logo.svg';

import styles from './index.module.scss';

const MilMoveHeader = ({ children }) => (
  <Header basic className={styles.mmHeader}>
    <div className="usa-nav-container">
      <div className="usa-navbar">
        <Title>
          <a href="/" title="Home" aria-label="Home">
            <img src={MmLogo} alt="MilMove" />
          </a>
        </Title>
      </div>
      <div className={styles.contents}>{children}</div>
    </div>
  </Header>
);

MilMoveHeader.defaultProps = {
  children: null,
};

MilMoveHeader.propTypes = {
  children: node,
};

export default MilMoveHeader;
