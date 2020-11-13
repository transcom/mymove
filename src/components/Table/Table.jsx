/* eslint-disable react/no-array-index-key */
import React, { Fragment } from 'react';
import PropTypes from 'prop-types';
import { Button, Dropdown } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faChevronLeft, faChevronRight } from '@fortawesome/free-solid-svg-icons';

import styles from './Table.module.scss';

const Table = ({
  handleClick,
  showFilters,
  showPagination,
  getTableProps,
  getTableBodyProps,
  headerGroups,
  rows,
  prepareRow,
  canPreviousPage,
  canNextPage,
  pageSize,
  nextPage,
  previousPage,
  gotoPage,
  setPageSize,
  pageOptions,
  perPage,
  pageIndex,
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
              {showFilters && (
                <tr className={styles.tableHeaderFilters} key={`headerGroupFilters${hgIndex}`}>
                  {headerGroup.headers.map((column, headerIndex) => (
                    <th key={`headerFilter${headerIndex}`} data-testid={column.id}>
                      {/* isFilterable is a custom prop that can be set in the Column object */}
                      <div>{column.isFilterable ? column.render('Filter') : null}</div>
                    </th>
                  ))}
                </tr>
              )}
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
              id="table-rows-per-page"
              className={styles.paginationSelect}
              name="table-rows-per-page"
              defaultValue={pageSize}
              onChange={(e) => {
                setPageSize(Number(e.target.value));
              }}
            >
              {perPage.map((page, index) => (
                <option value={page} key={`page-size-${index}`}>
                  {page}
                </option>
              ))}
            </Dropdown>
            <div>rows per page</div>
          </div>
          <div className={styles.tableControlPagination}>
            <Button
              type="button"
              unstyled
              className={`${styles.pageControlButton} ${styles.pageControlButtonPrev}`}
              onClick={previousPage}
              disabled={!canPreviousPage}
            >
              <FontAwesomeIcon className={`${styles.paginationIconLeft} fas fa-chevron-left`} icon={faChevronLeft} />
              <span>Prev</span>
            </Button>
            <Dropdown
              id="table-pagination"
              className={styles.paginationSelect}
              name="table-pagination"
              value={pageIndex}
              onChange={(e) => gotoPage(Number(e.target.value))}
            >
              {pageOptions.length > 0 ? (
                pageOptions.map((pageOption, index) => (
                  <option value={pageOption} key={`page-options-${index}`}>
                    {pageOption + 1}
                  </option>
                ))
              ) : (
                <option value={0}>{1}</option>
              )}
            </Dropdown>
            <Button
              type="button"
              unstyled
              className={`${styles.pageControlButton} ${styles.pageControlButtonNext}`}
              onClick={nextPage}
              disabled={!canNextPage}
            >
              <span>Next</span>
              <FontAwesomeIcon className={`${styles.paginationIconRight} fas fa-chevron-right`} icon={faChevronRight} />
            </Button>
          </div>
        </div>
      )}
    </div>
  );
};

Table.propTypes = {
  handleClick: PropTypes.func,
  showFilters: PropTypes.bool,
  showPagination: PropTypes.bool,
  previousPage: PropTypes.func,
  nextPage: PropTypes.func,
  setPageSize: PropTypes.func,
  gotoPage: PropTypes.func,
  // below are props from useTable() hook
  getTableProps: PropTypes.func.isRequired,
  getTableBodyProps: PropTypes.func.isRequired,
  headerGroups: PropTypes.arrayOf(PropTypes.object).isRequired,
  rows: PropTypes.arrayOf(PropTypes.object).isRequired,
  prepareRow: PropTypes.func.isRequired,
  canPreviousPage: PropTypes.bool,
  canNextPage: PropTypes.bool,
  pageCount: PropTypes.number,
  pageIndex: PropTypes.number,
  pageSize: PropTypes.number,
  pageOptions: PropTypes.arrayOf(PropTypes.number),
  perPage: PropTypes.arrayOf(PropTypes.number),
};

Table.defaultProps = {
  handleClick: undefined,
  showFilters: false,
  showPagination: false,
  canPreviousPage: undefined,
  previousPage: undefined,
  nextPage: undefined,
  setPageSize: undefined,
  gotoPage: undefined,
  canNextPage: undefined,
  pageCount: undefined,
  pageIndex: 0,
  pageSize: 20,
  pageOptions: [0],
  perPage: [10, 20, 50],
};

export default Table;
