import React from 'react';
import { useTable, useFilters } from 'react-table';
import PropTypes from 'prop-types';

import { textFilter } from './utils';
import TextBoxFilter from './Filters/TextBoxFilter';
import styles from './Table.module.scss';

const Table = ({ data, columns, hiddenColumns, handleClick }) => {
  const filterTypes = React.useMemo(
    () => ({
      // "startWith"
      text: textFilter,
    }),
    [],
  );

  const defaultColumn = React.useMemo(
    () => ({
      // Let's set up our default Filter UI
      Filter: TextBoxFilter,
    }),
    [],
  );

  const tableData = React.useMemo(() => data, [data]);
  const tableColumns = React.useMemo(() => columns, [columns]);
  const { getTableProps, getTableBodyProps, headerGroups, rows, prepareRow } = useTable(
    {
      columns: tableColumns,
      data: tableData,
      initialState: { hiddenColumns },
      defaultColumn, // Be sure to pass the defaultColumn option
      filterTypes,
    },
    useFilters,
  );

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
                  <div>{column.canFilter ? column.render('Filter') : null}</div>
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
                {row.cells.map((cell) => {
                  return (
                    <td data-testid={`${cell.column.id}-${cell.row.id}`} {...cell.getCellProps()}>
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
      accessor: PropTypes.oneOfType([PropTypes.string, PropTypes.func]),
    }),
  ),
  hiddenColumns: PropTypes.arrayOf(PropTypes.string),
  handleClick: PropTypes.func,
};

Table.defaultProps = {
  data: [],
  columns: [],
  hiddenColumns: [],
  handleClick: undefined,
};

export default Table;
