import { useQuery } from '@tanstack/react-query';
import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

import TableQueue from 'components/Table/TableQueue';
import { createHeader } from 'components/Table/utils';
import { PRIME_SIMULATOR_AVAILABLE_MOVES } from 'constants/queryKeys';
// TODO: This is very clunky. There are shared/formatters and util/formatters
// that determine dates. This way is a way to do it now, but this should be
// refactored as part of TRA work to be done differently across the app.
// For now though, I'm going to be using the `formatDateFromIso` function and
// then leverage a constant for how the date should be formatted.
import { getPrimeSimulatorAvailableMoves } from 'services/primeApi';
import LoadingPlaceholder from 'shared/LoadingPlaceholder';
import SomethingWentWrong from 'shared/SomethingWentWrong';
import { DATE_TIME_FORMAT_STRING } from 'shared/constants';
import { getQueriesStatus } from 'utils/api';
import { formatDateFromIso } from 'utils/formatters';

const columnHeaders = () => [
  createHeader('Move ID', 'id', {
    id: 'id',
    isFilterable: true,
  }),
  createHeader('Move code', 'moveCode', {
    id: 'moveCode',
    isFilterable: true,
  }),
  createHeader(
    'Created at',
    (row) => {
      return formatDateFromIso(row.createdAt, DATE_TIME_FORMAT_STRING);
    },
    { id: 'createdAt' },
  ),
  createHeader(
    'Updated at',
    (row) => {
      return formatDateFromIso(row.updatedAt, DATE_TIME_FORMAT_STRING);
    },
    { id: 'updatedAt' },
  ),
  createHeader('e-Tag', 'eTag'),
  createHeader('Order ID', 'orderID'),
  createHeader('Type', 'ppmType'),
  createHeader('Reference ID', 'referenceId'),
  createHeader(
    'Available to Prime at',
    (row) => {
      return formatDateFromIso(row.availableToPrimeAt, DATE_TIME_FORMAT_STRING);
    },
    { id: 'availableToPrimeAt' },
  ),
];

const PrimeSimulatorAvailableMoves = () => {
  const navigate = useNavigate();

  const todayDate = new Date();
  const [dateSelected, setDateSelected] = useState(todayDate.toISOString().split('T')[0]);
  const { data = {}, ...primeSimulatorAvailableMovesQuery } = useQuery(
    [PRIME_SIMULATOR_AVAILABLE_MOVES, { date: `${dateSelected}` }],
    ({ queryKey: [key, { ...date }] }) => {
      return getPrimeSimulatorAvailableMoves(key, date);
    },
  );

  const apiQuery = () => {
    const { isLoading, isError, isSuccess } = getQueriesStatus([primeSimulatorAvailableMovesQuery]);
    // README: This queueResult is being artificially constructed rather than
    // created using the `..dataProp` destructering of other functions because
    // the Prime API does not return an Object that the TableQueue component can
    // consume. So the queueResult mimics that Objects properties since `data` in
    // this case is a simple Array of Prime Available Moves.
    const queueResult = {
      data,
      page: 1,
      perPage: data.length,
      totalCount: data.length,
    };
    return {
      queueResult,
      isLoading,
      isError,
      isSuccess,
    };
  };
  const { isLoading, isError } = apiQuery();
  if (isLoading) return <LoadingPlaceholder />;
  if (isError) return <SomethingWentWrong />;

  const setFilterByDate = () => {
    const filterDate = document.getElementById('filterDate').value;
    const dateRegex = /^\d{4}-(0[1-9]|1[012])-(0[1-9]|[12][0-9]|3[01])$/;

    if (dateRegex.exec(filterDate)) {
      setDateSelected(filterDate);
      apiQuery();

      document.getElementById('error').innerHTML = '&nbsp;';
    } else {
      document.getElementById('error').innerHTML = 'Enter a valid date.';
    }
  };

  return (
    <>
      <div>
        <p>Select Filter from Date: (YYYY-MM-DD)</p>
        <p id="error">&nbsp;</p>
        <input type="text" id="filterDate" defaultValue={dateSelected} data-testid="prime-date-filter" />
        <button type="button" onClick={setFilterByDate}>
          Filter
        </button>
      </div>
      <TableQueue
        title="Moves available to Prime"
        columns={columnHeaders()}
        useQueries={apiQuery}
        handleClick={(row) => {
          navigate(`/simulator/moves/${row.id}/details`);
        }}
        defaultSortedColumns={[{ id: 'availableToPrimeAt', desc: false }]}
        defaultHiddenColumns={['eTag']}
        defaultCanSort
        disableSortBy={false}
        showFilters
        showPagination
        manualFilters={false}
      />
    </>
  );
};

export default PrimeSimulatorAvailableMoves;
