import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faTruck } from '@fortawesome/free-solid-svg-icons';

import styles from './MultiMovesMoveHeader.module.scss';

const MultiMovesMoveHeader = ({ title }) => {
  return (
    <div className={styles.moveHeaderContainer}>
      <FontAwesomeIcon icon={faTruck} data-testid="truck-icon" />
      <h3>{title}</h3>
    </div>
  );
};

export default MultiMovesMoveHeader;
