import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './Helper.module.scss';

const Helper = ({ children, title, className }) => (
  <div className={classnames(styles.Helper, className)}>
    <h3 className={styles.header}>{title}</h3>
    {children}
  </div>
);

Helper.propTypes = {
  title: PropTypes.string.isRequired,
  children: PropTypes.node,
  className: PropTypes.string,
};

Helper.defaultProps = {
  children: null,
  className: '',
};

export default Helper;
