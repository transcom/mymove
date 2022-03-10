import React from 'react';
import { string } from 'prop-types';

import TableQueue from 'components/Table/TableQueue';
import { createHeader } from 'components/Table/utils';
import { useGHCGetMoveHistory } from 'hooks/queries';
import { formatDateFromIso } from 'shared/formatters';

const formatChangedValues = (changedValues) => {
  return changedValues
    ? changedValues.map((changedValue) => `${changedValue.columnName}: ${changedValue.columnValue}`).join(', ')
    : '';
};

const columns = [
  createHeader('Date & Time', (row) => formatDateFromIso(`${row.actionTstampClk}`, 'DD MMM YY HH:mm')),
  createHeader('Event', 'eventName'),
  createHeader('Details', (row) => formatChangedValues(row.changedValues)),
  createHeader('Modified By', 'user.name'),
];

const MoveHistory = ({ moveCode }) => {
  const useGetMoveHistoryQuery = () => {
    return useGHCGetMoveHistory(moveCode);
  };

  return (
    <TableQueue
      showFilters={false}
      showPagination={false}
      disableSortBy
      columns={columns}
      title="Move history"
      handleClick={() => {}}
      useQueries={useGetMoveHistoryQuery}
    />
  );
};

MoveHistory.propTypes = {
  moveCode: string.isRequired,
};

export default MoveHistory;
