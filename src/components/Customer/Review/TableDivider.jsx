import React from 'react';
import { string } from 'prop-types';

const TableDivider = ({ className }) => (
  <tr>
    <td
      className={className}
      colSpan="100%"
      style={{ paddingTop: 0, paddingBottom: 0, borderTop: 'none', borderBottom: '1px solid black' }}
    />
  </tr>
);

TableDivider.propTypes = {
  className: string,
};

TableDivider.defaultProps = {
  className: '',
};

export default TableDivider;
