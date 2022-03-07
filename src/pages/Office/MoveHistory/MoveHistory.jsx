import React from 'react';
import { string } from 'prop-types';

import TableQueue from 'components/Table/TableQueue';
import { createHeader } from 'components/Table/utils';
import { useGHCGetMoveHistory } from 'hooks/queries';

const formatDetails = (changedValue) => {
  return `${changedValue.columnName}: ${changedValue.columnValue}`;
};

const columns = [
  createHeader('Date & Time', 'actionTstampClk'),
  createHeader('Event', 'eventName'),
  // createHeader('Details', (row) => `${row.changedValues.map((changedValue) => changedValue.columnName)}`),
  createHeader('Details', (row) => `${row.changedValues.map((changedValue) => formatDetails(changedValue))}`),
  createHeader('User', 'user.name'),
];

const handleClick = () => {};

const MoveHistory = ({ moveCode }) => {
  // const { moveHistory, isLoading, isError } = useGHCGetMoveHistory(moveCode);

  // if (isLoading) return <LoadingPlaceholder />;
  // if (isError) return <SomethingWentWrong />;

  const useGetMoveHistoryQuery = ({ sort, order, currentPage, currentPageSize }) => {
    return useGHCGetMoveHistory({ moveCode, sort, order, currentPage, currentPageSize });
  };

  return (
    <TableQueue
      showFilters={false}
      showPagination
      manualSortBy
      defaultCanSort
      defaultSortedColumns={[{ id: 'actionTstampClk', desc: true }]}
      disableMultiSort
      disableSortBy={false}
      columns={columns}
      title="Move history"
      handleClick={handleClick}
      useQueries={useGetMoveHistoryQuery}
    />
  );
};

MoveHistory.propTypes = {
  moveCode: string.isRequired,
};

export default MoveHistory;
