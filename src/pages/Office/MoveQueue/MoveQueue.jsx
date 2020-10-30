import React, { useState, useEffect, useMemo } from 'react';
import { withRouter } from 'react-router-dom';
import { GridContainer } from '@trussworks/react-uswds';
import { useTable, useFilters } from 'react-table';

import styles from './MoveQueue.module.scss';

import { HistoryShape } from 'types/router';
import Table from 'components/Table/Table';
import { createHeader } from 'components/Table/utils';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { useMovesQueueQueries } from 'hooks/queries';
import { serviceMemberAgencyLabel } from 'shared/formatters';
import TextBoxFilter from 'components/Table/Filters/TextBoxFilter';
import MultiSelectCheckBoxFilter from 'components/Table/Filters/MultiSelectCheckBoxFilter';
import SelectFilter from 'components/Table/Filters/SelectFilter';
import { MOVE_STATUS_OPTIONS } from 'constants/queues';

const moveStatusOptions = Object.keys(MOVE_STATUS_OPTIONS).map((key) => ({
  value: key,
  label: MOVE_STATUS_OPTIONS[`${key}`],
}));

const branchFilterOptions = [
  { value: '', label: 'All' },
  { value: 'ARMY', label: 'Army' },
  { value: 'NAVY', label: 'Navy' },
  { value: 'MARINES', label: 'Marine Corps' },
  { value: 'AIR_FORCE', label: 'Air Force' },
  { value: 'COAST_GUARD', label: 'Coast Guard' },
];

const columns = [
  createHeader('ID', 'id'),
  createHeader(
    'Customer name',
    (row) => {
      return `${row.customer.last_name}, ${row.customer.first_name}`;
    },
    {
      id: 'lastName',
      isFilterable: true,
    },
  ),
  createHeader('DoD ID', 'customer.dodID', {
    id: 'dodID',
    isFilterable: true,
  }),
  createHeader('Status', 'status', {
    isFilterable: true,
    // eslint-disable-next-line react/jsx-props-no-spreading
    Filter: (props) => <MultiSelectCheckBoxFilter options={moveStatusOptions} {...props} />,
  }),
  createHeader('Move Code', 'locator', {
    id: 'moveID',
    isFilterable: true,
  }),
  createHeader(
    'Branch',
    (row) => {
      return serviceMemberAgencyLabel(row.customer.agency);
    },
    {
      id: 'branch',
      isFilterable: true,
      // eslint-disable-next-line react/jsx-props-no-spreading
      Filter: (props) => <SelectFilter options={branchFilterOptions} {...props} />,
    },
  ),
  createHeader('# of shipments', 'shipmentsCount'),
  createHeader('Destination duty station', 'destinationDutyStation.name', {
    id: 'destinationDutyStation',
    isFilterable: true,
  }),
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
      manualFilters: true,
    },
    useFilters,
  );

  // When these table states change, fetch new data!
  useEffect(() => {
    if (!isLoading && !isError) {
      setParamFilters(filters);
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
