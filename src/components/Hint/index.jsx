import React from 'react';
import { string, node, bool } from 'prop-types';

import styles from './index.module.scss';

const Hint = ({ className, children, darkText, ...props }) => (
  // eslint-disable-next-line react/jsx-props-no-spreading
  <div {...props} className={`${styles.Hint} ${className} ${darkText ? styles.DarkText : ''}`}>
    {children}
  </div>
);

Hint.propTypes = {
  className: string,
  children: node.isRequired,
  darkText: bool,
};

Hint.defaultProps = {
  className: '',
  darkText: false,
};

export default Hint;
