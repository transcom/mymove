import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './index.module.scss';

const DataTable = ({ columnHeaders, dataRow, icon, custClass }) => (
  <table className={classnames(styles.dataTable, 'table--data-point', custClass)}>
    <thead className="table--small">
      <tr>
        {columnHeaders.map((header) => (
          <th key={header}>{header}</th>
        ))}
      </tr>
    </thead>
    <tbody>
      <tr>
        {dataRow.map((cell, i) => (
          // RA Summary: eslint:react/no-array-index-key
          // RA: Using the index as an element key in cases where the array is reordered will result in unnecessary renders.
          // RA: The source data is unstructured, with a potential for duplicate values amongst siblings.
          // RA: A reorder function is not implemented for this array.
          // RA Developer Status: Mitigated
          // RA Validator Status: Mitigated
          // RA Modified Severity: N/A
          // no unique identifier that can be used as a key, cell values can be duplicates (e.g. Dates)
          // eslint-disable-next-line react/no-array-index-key
          <td key={i}>
            <div className={classnames({ [`${styles.iconCellContainer}`]: !!icon && i === 0 })}>
              <span>{cell}</span>
              {!!icon && i === 0 && icon}
            </div>
          </td>
        ))}
      </tr>
    </tbody>
  </table>
);

DataTable.propTypes = {
  columnHeaders: PropTypes.arrayOf(PropTypes.node).isRequired,
  dataRow: PropTypes.arrayOf(PropTypes.node).isRequired,
  icon: PropTypes.node,
  custClass: PropTypes.string,
};

DataTable.defaultProps = {
  icon: null,
  custClass: '',
};

export default DataTable;
