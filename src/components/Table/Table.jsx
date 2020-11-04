/* eslint-disable react/no-array-index-key */
import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { Button, Dropdown } from '@trussworks/react-uswds';

import { ReactComponent as ChevronLeft } from '../../shared/icon/chevron-left.svg';
import { ReactComponent as ChevronRight } from '../../shared/icon/chevron-right.svg';

import styles from './Table.module.scss';

const Table = ({ handleClick, getTableProps, getTableBodyProps, headerGroups, rows, prepareRow }) => {
  return (
    /* eslint-disable react/jsx-props-no-spreading */
    <div className={styles.inlineBlock}>
      <div data-testid="react-table" className={styles.Table}>
        <table {...getTableProps()}>
          <thead>
            {headerGroups.map((headerGroup, hgIndex) => (
              <Fragment key={`headerGroup${hgIndex}`}>
                <tr {...headerGroup.getHeaderGroupProps()}>
                  {headerGroup.headers.map((column, headerIndex) => (
                    <th key={`header${headerIndex}`} data-testid={column.id} {...column.getHeaderProps()}>
                      {column.render('Header')}
                    </th>
                  ))}
                </tr>
                <tr className={styles.tableHeaderFilters} key={`headerGroupFilters${hgIndex}`}>
                  {headerGroup.headers.map((column, headerIndex) => (
                    <th key={`headerFilter${headerIndex}`} data-testid={column.id}>
                      {/* isFilterable is a custom prop that can be set in the Column object */}
                      <div>{column.isFilterable ? column.render('Filter') : null}</div>
                    </th>
                  ))}
                </tr>
              </Fragment>
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
                      <td
                        key={`cell${index}`}
                        data-testid={`${cell.column.id}-${cell.row.id}`}
                        {...cell.getCellProps()}
                      >
                        {cell.render('Cell')}
                      </td>
                    );
                  })}
                </tr>
              );
            })}
          </tbody>
        </table>
        <div className={styles.paginationSectionWrapper}>
          <div className={styles.paginationTableWrapper}>
            <div className={styles.tableControlRowsPerPage}>
              <Dropdown className={styles.usaSelect} name="table-rows-per-page">
                <option value="10">10</option>
                <option value="20">20</option>
                <option value="50">50</option>
              </Dropdown>
              <p>rows per page</p>
            </div>
          </div>
          <div className={styles.tableControlPagination}>
            <Button disabled className={styles.usaButtonUnstyled}>
              <span className="icon">
                <ChevronLeft />
              </span>
              <span>Prev</span>
            </Button>
            <Dropdown className={styles.usaSelect} name="table-pagination">
              <option value="1">1</option>
              <option value="2">2</option>
              <option value="3">3</option>
            </Dropdown>
            <Button className={styles.usaButtonUnstyled}>
              <span>Next</span>
              <span className="icon">
                <ChevronRight />
              </span>
            </Button>
          </div>
        </div>
      </div>
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
