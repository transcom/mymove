import React from 'react';
import PropTypes from 'prop-types';
import { Header, Title } from '@trussworks/react-uswds';

import MmLogo from '../../shared/images/milmove-logo.svg';

import styles from './index.module.scss';

const MilMoveHeader = ({ isSpecialMove, children }) => {
  return (
    <Header basic className={styles.mmHeader}>
      <div className="usa-nav-container">
        <div className="usa-navbar">
          <Title>
            <a href="/" title="Home" aria-label="Home">
              <img src={MmLogo} alt="MilMove" />
            </a>
          </Title>
        </div>
        {isSpecialMove ? (
          <div data-testid="specialMovesLabel" className={styles.specialMovesLabel}>
            <p>BLUEBARK</p>
          </div>
        ) : null}
        <div className={styles.contents}>{children}</div>
      </div>
    </Header>
  );
};

MilMoveHeader.defaultProps = {
  isSpecialMove: null,
  children: null,
};

MilMoveHeader.propTypes = {
  children: PropTypes.node,
  isSpecialMove: PropTypes.bool,
};

export default MilMoveHeader;
