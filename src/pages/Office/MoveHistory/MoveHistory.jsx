import React from 'react';
import { string } from 'prop-types';

import styles from './MoveHistory.module.scss';

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
  createHeader('User', 'user.name'),
];

const MoveHistory = ({ moveCode }) => {
  const useGetMoveHistoryQuery = ({ sort, order, currentPage, currentPageSize }) => {
    return useGHCGetMoveHistory({ moveCode, sort, order, currentPage, currentPageSize });
  };

  return (
    <div className={styles.MoveHistory}>
      <TableQueue
        showFilters={false}
        showPagination={false}
        defaultSortedColumns={[{ id: 'Date & Time', desc: true }]}
        disableSortBy
        columns={columns}
        title="Move history"
        handleClick={() => {}}
        useQueries={useGetMoveHistoryQuery}
      />
    </div>
  );
};

MoveHistory.propTypes = {
  moveCode: string.isRequired,
};

export default MoveHistory;
