import React from 'react';
import { useTable } from 'react-table';
import PropTypes from 'prop-types';

import styles from './Table.module.scss';

const Table = ({ data, columns }) => {
  const tableData = React.useMemo(() => data, [data]);
  const tableColumns = React.useMemo(() => columns, [columns]);
  const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } = useTable({
    columns: tableColumns,
    data: tableData,
  });

  return (
    /* eslint-disable react/jsx-props-no-spreading */
    <div data-testid="react-table" className={styles.Table}>
      <table {...getTableProps()}>
        <thead>
          {headerGroups.map((headerGroup) => (
            <tr {...headerGroup.getHeaderGroupProps()}>
              {headerGroup.headers.map((column) => (
                <th data-testid={column.id} {...column.getHeaderProps()}>
                  {column.render('Header')}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody {...getTableBodyProps()}>
          {rows.map((row) => {
            prepareRow(row);
            return (
              <tr {...row.getRowProps()}>
                {row.cells.map((cell) => {
                  return (
                    <td data-testid={cell.column.id} {...cell.getCellProps()}>
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
  // data is an array of objects to populate the table
  data: PropTypes.arrayOf(PropTypes.object),
  // columns is an array of objects to define the header name and accessor
  columns: PropTypes.arrayOf(
    PropTypes.shape({
      Header: PropTypes.string,
      accessor: PropTypes.string,
    }),
  ),
};

Table.defaultProps = {
  data: [],
  columns: [],
};

export default Table;
