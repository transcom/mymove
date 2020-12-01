import React from 'react';
import classnames from 'classnames';
import propTypes from 'prop-types';

import styles from './index.module.scss';

const DataPointGroup = ({ children, className }) => {
  return <div className={classnames(styles.dataPointGroup, 'table--data-point-group', className)}>{children}</div>;
};

DataPointGroup.propTypes = {
  className: propTypes.string,
  children: propTypes.node.isRequired,
};

DataPointGroup.defaultProps = {
  className: '',
};

export default DataPointGroup;
