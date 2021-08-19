/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './index.module.scss';

const Callout = ({ children, className, ...props }) => (
  <div className={classnames(styles.callOutContainer, className)} {...props}>
    {children}
  </div>
);

Callout.propTypes = {
  children: PropTypes.node.isRequired,
  className: PropTypes.string,
};

Callout.defaultProps = {
  className: '',
};

export default Callout;
