import React from 'react';
import { string } from 'prop-types';

import styles from './Review.module.scss';

const TableDivider = ({ className }) => (
  <tr>
    <td className={`${styles['table-divider']} ${className}`} colSpan="100%" />
  </tr>
);

TableDivider.propTypes = {
  className: string,
};

TableDivider.defaultProps = {
  className: '',
};

export default TableDivider;
