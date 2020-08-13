import React from 'react';
import classnames from 'classnames';
import propTypes from 'prop-types';

import styles from './index.module.scss';

const DataPointGroup = ({ children }) => {
  return <div className={classnames(styles.dataPointGroup, 'table--data-point-group')}>{children}</div>;
};

DataPointGroup.propTypes = {
  children: propTypes.node.isRequired,
};

export default DataPointGroup;
