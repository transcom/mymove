import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './LeftNav.module.scss';

const LeftNav = ({ className, children }) => <nav className={classnames(styles.LeftNav, className)}>{children}</nav>;

LeftNav.propTypes = {
  className: PropTypes.string,
  children: PropTypes.node.isRequired,
};

LeftNav.defaultProps = {
  className: '',
};

export default LeftNav;
