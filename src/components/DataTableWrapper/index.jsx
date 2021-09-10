import React from 'react';
import classnames from 'classnames';
import propTypes from 'prop-types';

import styles from './index.module.scss';

const DataTableWrapper = ({ children, className }) => {
  return <div className={classnames(styles.DataTableWrapper, 'table--data-point-group', className)}>{children}</div>;
};

DataTableWrapper.propTypes = {
  className: propTypes.string,
  children: propTypes.node.isRequired,
};

DataTableWrapper.defaultProps = {
  className: '',
};

export default DataTableWrapper;
