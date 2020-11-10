/* eslint-disable react/no-array-index-key */
import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { Button, Dropdown } from '@trussworks/react-uswds';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import { faChevronLeft, faChevronRight } from '@fortawesome/fontawesome-free-solid/';

import styles from './Table.module.scss';

const Table = ({
  handleClick,
  getTableProps,
  getTableBodyProps,
  headerGroups,
  rows,
  prepareRow,
  canPreviousPage,
  canNextPage,
  showPagination,
  pageSize,
  nextPage,
  handlePreviousClick,
  handlePageSelect,
  setPageSize,
  pageIndex,
  pageOptions,
}) => {
  return (
    /* eslint-disable react/jsx-props-no-spreading */
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
      {showPagination && (
        <div className={styles.paginationSectionWrapper} data-testid="pagination">
          <div className={styles.tableControlRowsPerPage}>
            <Dropdown
              className={styles.usaSelect}
              name="table-rows-per-page"
              defaultValue={pageSize}
              onChange={(e) => {
                setPageSize(Number(e.target.value));
              }}
            >
              {pageOptions.map((size, index) => (
                <option value={size} key={`page-size-${index}`}>
                  {size}
                </option>
              ))}
            </Dropdown>
            <div>rows per page</div>
          </div>
          <div className={styles.tableControlPagination}>
            <Button
              className={styles.usaButtonUnstyled}
              onClick={() => handlePreviousClick(pageIndex)}
              disabled={!canPreviousPage}
            >
              <FontAwesomeIcon className="icon fas fa-chevron-left" icon={faChevronLeft} />
              <span>Prev</span>
            </Button>
            <Dropdown
              className={styles.usaSelect}
              name="table-pagination"
              onChange={(e) => {
                handlePageSelect(e);
              }}
            >
              <option value="1">1</option>
              <option value="2">2</option>
              <option value="3">3</option>
            </Dropdown>
            <Button className={styles.usaButtonUnstyled} onClick={nextPage}>
              <span>Next</span>
              <FontAwesomeIcon className="icon fas fa-chevron-right" icon={faChevronRight} />
            </Button>
          </div>
        </div>
      )}
    </div>
  );
};

Table.propTypes = {
  handleClick: PropTypes.func,
  handlePreviousClick: PropTypes.func,
  nextPage: PropTypes.func,
  setPageSize: PropTypes.func,
  handlePageSelect: PropTypes.func,
  // below are props from useTable() hook
  getTableProps: PropTypes.func.isRequired,
  getTableBodyProps: PropTypes.func.isRequired,
  headerGroups: PropTypes.arrayOf(PropTypes.object).isRequired,
  rows: PropTypes.arrayOf(PropTypes.object).isRequired,
  prepareRow: PropTypes.func.isRequired,
  showPagination: PropTypes.bool,
  canPreviousPage: PropTypes.bool,
  canNextPage: PropTypes.bool,
  pageCount: PropTypes.number,
  pageIndex: PropTypes.number,
  pageSize: PropTypes.number,
  state: PropTypes.node,
  pageOptions: PropTypes.arrayOf(PropTypes.number),
};

Table.defaultProps = {
  handleClick: undefined,
  showPagination: false,
  pageIndex: 0,
  pageSize: 20,
  canPreviousPage: false,
  pageOptions: [10, 20, 50],
};

export default Table;
