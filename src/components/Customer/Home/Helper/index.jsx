import React from 'react';
import { string, object, node } from 'prop-types';

import styles from './Helper.module.scss';

const Helper = ({ containerStyles, children, title }) => (
  <div className={styles['helper-container']} style={containerStyles}>
    <h3 className={styles['helper-header']}>{title}</h3>
    {children}
  </div>
);

Helper.propTypes = {
  // eslint-disable-next-line react/forbid-prop-types
  containerStyles: object,
  title: string.isRequired,
  children: node,
};

Helper.defaultProps = {
  containerStyles: {},
  children: null,
};

export default Helper;
