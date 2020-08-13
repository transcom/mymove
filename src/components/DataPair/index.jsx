import React from 'react';
import classnames from 'classnames';
import propTypes from 'prop-types';

import styles from './index.module.scss';

const DataPair = ({ children }) => {
  return <div className={classnames(styles.dataPair, 'table--data-pair')}>{children}</div>;
};

DataPair.propTypes = {
  children: propTypes.node.isRequired,
};

export default DataPair;
