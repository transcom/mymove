import React from 'react';

import styles from './StyledLine.module.scss';

export const StyledLine = ({ width, color, className }) => {
  return (
    <div
      className={`${className || styles.styledLine}`}
      style={{
        width: width || '75%',
        backgroundColor: color || '#565c65',
      }}
    />
  );
};

export default StyledLine;
