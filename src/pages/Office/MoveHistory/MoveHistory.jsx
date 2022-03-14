import React from 'react';
import { string } from 'prop-types';

import styles from './MoveHistory.module.scss';

import TableQueue from 'components/Table/TableQueue';
import { createHeader } from 'components/Table/utils';
import { useGHCGetMoveHistory } from 'hooks/queries';
import { formatDateFromIso } from 'shared/formatters';

const formatChangedValues = (changedValues) => {
  return changedValues
    ? changedValues.map((changedValue) => (
        <div key={`${changedValue.columnName}-${changedValue.columnValue}`}>
          {changedValue.columnName}: {changedValue.columnValue}
        </div>
      ))
    : '';
};

const columns = [
  createHeader(
    'Date & Time',
    (row) => <div className={styles.dateAndTime}>{formatDateFromIso(row.actionTstampClk, 'DD MMM YY HH:mm')}</div>,
    { id: 'move-history-date-time' },
  ),
  createHeader('Event', (row) => <div className={styles.event}>{row.eventName}</div>, {
    id: 'move-history-event',
  }),
  createHeader(
    'Details',
    (row) => {
      return <div className={styles.details}>{formatChangedValues(row.changedValues)}</div>;
    },
    { id: 'move-history-details' },
  ),
  createHeader('Modified By', (row) => <div className={styles.modifiedBy}>{row.userName}</div>, {
    id: 'move-history-modified-by',
  }),
];

const MoveHistory = ({ moveCode }) => {
  const useGetMoveHistoryQuery = () => {
    return useGHCGetMoveHistory(moveCode);
  };

  return (
    <div className={styles.MoveHistoryTable}>
      <TableQueue
        showFilters={false}
        showPagination={false}
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
