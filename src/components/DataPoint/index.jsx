import React from 'react';
import PropTypes from 'prop-types';
import classnames from 'classnames';

import styles from './index.module.scss';

const DataPoint = ({ columnHeaders, dataRow, icon, custClass }) => (
  <table className={classnames(styles.dataPoint, 'table--data-point', custClass)}>
    <thead className="table--small">
      <tr>
        {columnHeaders.map((header, i) => (
          // eslint-disable-next-line react/no-array-index-key
          <th key={i}>{header}</th>
        ))}
      </tr>
    </thead>
    <tbody>
      <tr>
        {dataRow.map((cell, i) => (
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
