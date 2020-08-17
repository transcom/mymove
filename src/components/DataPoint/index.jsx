import React from 'react';
import propTypes from 'prop-types';
import classnames from 'classnames';

import styles from './index.module.scss';

const DataPoint = ({ columnHeaders, dataRow, Icon, custClass }) => (
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
        {dataRow.map((cell, i) => (
          <td key={columnHeaders[`${i}`]}>
            <div className={classnames({ [`${styles.iconCellContainer}`]: !!Icon && i === 0 })}>
              <span>{cell}</span>
              {!!Icon && i === 0 && <Icon />}
            </div>
          </td>
        ))}
      </tr>
    </tbody>
  </table>
);

DataPoint.propTypes = {
  columnHeaders: propTypes.arrayOf(propTypes.node).isRequired,
  dataRow: propTypes.arrayOf(propTypes.node).isRequired,
  Icon: propTypes.elementType,
  custClass: propTypes.string,
};

DataPoint.defaultProps = {
  Icon: undefined,
  custClass: '',
};

export default DataPoint;
