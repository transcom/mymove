import React from 'react';
import { string } from 'prop-types';

import styles from './Review.module.scss';

export const TableDivider = ({ className }) => (
  <tr>
    <td className={`${styles['table-divider']} ${className}`} data-testid="tableDivider" colSpan="100%" />
  </tr>
);

TableDivider.propTypes = {
  className: string,
};

TableDivider.defaultProps = {
  className: '',
};

export default TableDivider;
