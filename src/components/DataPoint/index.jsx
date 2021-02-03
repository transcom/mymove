import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './index.module.scss';

const DataPoint = ({ columnHeaders, dataRow, icon, custClass }) => (
  <table className={classnames(styles.dataPoint, 'table--data-point', custClass)}>
    <thead className="table--small">
      <tr>
        {columnHeaders.map((header) => (
          <th key={header}>{header}</th>
        ))}
      </tr>
    </thead>
    <tbody>
      <tr>
        {/* eslint-disable react/no-array-index-key */}
        {/* no unique identifier that can be used as a key, cell values can be duplicates (e.g. Dates) */}
        {dataRow.map((cell, i) => (
          <td key={i}>
            <div className={classnames({ [`${styles.iconCellContainer}`]: !!icon && i === 0 })}>
              <span>{cell}</span>
              {!!icon && i === 0 && icon}
            </div>
          </td>
        ))}
        {/* eslint-enable react/no-array-index-key */}
      </tr>
    </tbody>
  </table>
);

DataPoint.propTypes = {
  columnHeaders: PropTypes.arrayOf(PropTypes.node).isRequired,
  dataRow: PropTypes.arrayOf(PropTypes.node).isRequired,
  icon: PropTypes.node,
  custClass: PropTypes.string,
};

DataPoint.defaultProps = {
  icon: null,
  custClass: '',
};

export default DataPoint;
