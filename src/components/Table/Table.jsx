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
  initialState,
  prepareRow,
  canPreviousPage,
  canNextPage,
  gotoPage,
  nextPage,
  previousPage,
  setPageSize,
  showPagination,
  pageSize,
  handleNextClick,
  handlePageSelect,
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
              onClick={() => setPageSize()}
            >
              <option value="10">10</option>
              <option value="20">20</option>
              <option value="50">50</option>
            </Dropdown>
            <div>rows per page</div>
          </div>
          <div className={styles.tableControlPagination}>
            <Button
              className={styles.usaButtonUnstyled}
              onClick={() => previousPage(pageIndex)}
              disabled={!canPreviousPage}
            >
              <FontAwesomeIcon className="icon fas fa-chevron-left" icon={faChevronLeft} />
              <span>Prev</span>
            </Button>
            <Dropdown className={styles.usaSelect} name="table-pagination" onClick={() => handlePageSelect(rows.value)}>
              <option value="1">1</option>
              <option value="2">2</option>
              <option value="3">3</option>
            </Dropdown>
            <Button className={styles.usaButtonUnstyled} onClick={() => handleNextClick(pageIndex)}>
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
  handleNextClick: PropTypes.func,
  handlePageSelect: PropTypes.func,
  // below are props from useTable() hook
  getTableProps: PropTypes.func.isRequired,
  getTableBodyProps: PropTypes.func.isRequired,
  headerGroups: PropTypes.arrayOf(PropTypes.object).isRequired,
  rows: PropTypes.arrayOf(PropTypes.object).isRequired,
  prepareRow: PropTypes.func.isRequired,
  showPagination: PropTypes.bool,
  // eslint-disable-next-line react/forbid-prop-types
  initialState: PropTypes.object,
  canPreviousPage: PropTypes.bool,
  canNextPage: PropTypes.bool,
  pageCount: PropTypes.number,
  gotoPage: PropTypes.func,
  nextPage: PropTypes.func,
  previousPage: PropTypes.func,
  setPageSize: PropTypes.func,
  pageIndex: PropTypes.number,
  pageSize: PropTypes.number,
  state: PropTypes.node,
  goToPage: PropTypes.func,
};

Table.defaultProps = {
  handleClick: undefined,
  showPagination: false,
  initialState: { pageIndex: 0, pageSize: 20 },
  canPreviousPage: false,
};

export default Table;
