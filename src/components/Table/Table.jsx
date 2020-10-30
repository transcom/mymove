import React from 'react';
import PropTypes from 'prop-types';

import styles from './Table.module.scss';

const Table = ({ handleClick, getTableProps, getTableBodyProps, headerGroups, rows, prepareRow }) => {
  return (
    /* eslint-disable react/jsx-props-no-spreading */
    <div data-testid="react-table" className={styles.Table}>
      <table {...getTableProps()}>
        <thead>
          {headerGroups.map((headerGroup, hgIndex) => (
            // eslint-disable-next-line react/no-array-index-key
            <tr key={`headerGroup${hgIndex}`} {...headerGroup.getHeaderGroupProps()}>
              {headerGroup.headers.map((column, headerIndex) => (
                // eslint-disable-next-line react/no-array-index-key
                <th key={`header${headerIndex}`} data-testid={column.id} {...column.getHeaderProps()}>
                  {column.render('Header')}
                  {/* isFilterable is a custom prop that can be set in the Column object */}
                  {/* <div>{column.isFilterable ? column.render('Filter') : null}</div> */}
                  <div>{column.isFilterable ? column.render('Filter') : null}</div>
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody {...getTableBodyProps()}>
          {rows.map((row) => {
            prepareRow(row);
            return (
              <tr data-uuid={row.values.id} onClick={() => handleClick(row.values)} {...row.getRowProps()}>
                {row.cells.map((cell, index) => {
                  return (
                    // eslint-disable-next-line react/no-array-index-key
                    <td key={`cell${index}`} data-testid={`${cell.column.id}-${cell.row.id}`} {...cell.getCellProps()}>
                      {cell.render('Cell')}
                    </td>
                  );
                })}
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
};

Table.propTypes = {
  handleClick: PropTypes.func,
  // below are props from useTable() hook
  getTableProps: PropTypes.func.isRequired,
  getTableBodyProps: PropTypes.func.isRequired,
  headerGroups: PropTypes.arrayOf(PropTypes.object).isRequired,
  rows: PropTypes.arrayOf(PropTypes.object).isRequired,
  prepareRow: PropTypes.func.isRequired,
};

Table.defaultProps = {
  handleClick: undefined,
};

export default Table;
