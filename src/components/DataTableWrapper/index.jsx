import React from 'react';
import classnames from 'classnames';
import propTypes from 'prop-types';

import styles from './index.module.scss';

const DataTableWrapper = ({ children, className, testID }) => {
  return (
    <div className={classnames(styles.dataTableWrapper, className)} data-testid={testID}>
      {children}
    </div>
  );
};

DataTableWrapper.propTypes = {
  className: propTypes.string,
  children: propTypes.node.isRequired,
  testID: propTypes.string,
};

DataTableWrapper.defaultProps = {
  className: '',
  testID: '',
};

export default DataTableWrapper;
