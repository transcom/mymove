import React from 'react';
import { string } from 'prop-types';

import styles from './MoveHistory.module.scss';
import ModifiedBy from './ModifiedBy';
import MoveHistoryDetailsSelector from './MoveHistoryDetailsSelector';

import TableQueue from 'components/Table/TableQueue';
import { createHeader } from 'components/Table/utils';
import { useGHCGetMoveHistory } from 'hooks/queries';
import { formatDateFromIso } from 'shared/formatters';
import { getHistoryLogEventNameDisplay } from 'constants/historyLogUIDisplayName';

const columns = [
  createHeader(
    'Date & Time',
    (row) => <div className={styles.dateAndTime}>{formatDateFromIso(row.actionTstampClk, 'DD MMM YY HH:mm')}</div>,
    { id: 'move-history-date-time' },
  ),
  createHeader(
    'Event',
    (row) => (
      <div className={styles.event}>
        {getHistoryLogEventNameDisplay({ eventName: row.eventName, changedValues: row.changedValues })}
      </div>
    ),
    { id: 'move-history-event' },
  ),
  createHeader(
    'Details',
    (row) => {
      return (
        <MoveHistoryDetailsSelector
          eventName={row.eventName}
          oldValues={row.oldValues}
          changedValues={row.changedValues}
        />
      );
    },
    { id: 'move-history-details' },
  ),
  createHeader(
    'Modified By',
    (row) => (
      <ModifiedBy
        firstName={row.sessionUserFirstName}
        lastName={row.sessionUserLastName}
        email={row.sessionUserEmail}
        phone={row.sessionUserTelephone}
      />
    ),
    { id: 'move-history-modified-by' },
  ),
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
