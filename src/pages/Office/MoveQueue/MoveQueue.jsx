import React, { useState, useEffect, useMemo } from 'react';
import { withRouter } from 'react-router-dom';
import { GridContainer } from '@trussworks/react-uswds';
import { useTable, useFilters } from 'react-table';

import styles from './MoveQueue.module.scss';

import { HistoryShape } from 'types/router';
import Table from 'components/Table/Table';
import { createHeader, textFilter } from 'components/Table/utils';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useMovesQueueQueries } from 'hooks/queries';
import { departmentIndicatorLabel } from 'shared/formatters';
import TextBoxFilter from 'components/Table/Filters/TextBoxFilter';

const columns = [
  createHeader('ID', 'id'),
  createHeader(
    'Customer name',
    (row) => {
      return `${row.customer.last_name}, ${row.customer.first_name}`;
    },
    { id: 'lastName' },
    { isFilterable: true },
  ),
  createHeader('DoD ID', 'customer.dodID', { id: 'dodID' }, { isFilterable: true }),
  createHeader('Status', 'status', { isFilterable: true }),
  createHeader('Move ID', 'locator', { id: 'moveID' }, { isFilterable: true }),
  createHeader(
    'Branch',
    (row) => {
      return departmentIndicatorLabel(row.departmentIndicator);
    },
    { id: 'branch' },
  ),
  createHeader('# of shipments', 'shipmentsCount'),
  createHeader(
    'Destination duty station',
    'destinationDutyStation.name',
    { id: 'destinationDutyStation' },
    { isFilterable: true },
  ),
  createHeader('Origin GBLOC', 'originGBLOC'),
];

const MoveQueue = ({ history }) => {
  const [paramFilters, setParamFilters] = useState([]);

  const {
    queueMovesResult: { totalCount, queueMoves = [] },
    isLoading,
    isError,
  } = useMovesQueueQueries(paramFilters);

  // react-table setup below

  const filterTypes = useMemo(
    () => ({
      // "startWith"
      text: textFilter,
    }),
    [],
  );
  const defaultColumn = useMemo(
    () => ({
      // Let's set up our default Filter UI
      Filter: TextBoxFilter,
    }),
    [],
  );
  const tableData = useMemo(() => queueMoves, [queueMoves]);
  const tableColumns = useMemo(() => columns, []);
  const {
    getTableProps,
    getTableBodyProps,
    headerGroups,
    rows,
    prepareRow,
    state: { filters },
  } = useTable(
    {
      columns: tableColumns,
      data: tableData,
      initialState: { hiddenColumns: ['id'] },
      defaultColumn, // Be sure to pass the defaultColumn option
      filterTypes,
      manualFilters: true,
    },
    useFilters,
  );

  // When these table states change, fetch new data!
  useEffect(() => {
    if (!isLoading && !isError) {
      if (filters.length > 0) {
        setParamFilters(filters);
      }
    }
  }, [filters, isLoading, isError]);

  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const handleClick = (values) => {
    history.push(`/moves/${values.id}/details`);
  };

  return (
    <GridContainer containerSize="widescreen" className={styles.MoveQueue}>
      <h1>{`All moves (${totalCount})`}</h1>
      <div className={styles.tableContainer}>
        <Table
          handleClick={handleClick}
          getTableProps={getTableProps}
          getTableBodyProps={getTableBodyProps}
          headerGroups={headerGroups}
          rows={rows}
          prepareRow={prepareRow}
        />
      </div>
    </GridContainer>
  );
};

MoveQueue.propTypes = {
  history: HistoryShape.isRequired,
};

export default withRouter(MoveQueue);
