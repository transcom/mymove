import React from 'react';

import styles from './StyledLine.module.scss';

const StyledLine = ({ width, color }) => {
  return (
    <div
      className={styles.styledLine}
      style={{
        width: width || '75%',
        backgroundColor: color || '#565c65',
      }}
    />
  );
};

export default StyledLine;
