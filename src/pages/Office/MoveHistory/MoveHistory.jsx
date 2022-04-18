import React from 'react';
import { string } from 'prop-types';

import styles from './MoveHistory.module.scss';
import ModifiedBy from './ModifiedBy';
import MoveHistoryDetailsSelector from './MoveHistoryDetailsSelector';

import TableQueue from 'components/Table/TableQueue';
import { createHeader } from 'components/Table/utils';
import { useGHCGetMoveHistory } from 'hooks/queries';
import { formatDateFromIso } from 'shared/formatters';
import getMoveHistoryEventTemplate from 'constants/moveHistoryEventTemplate';

const columns = [
  createHeader(
    'Date & Time',
    (row) => <div className={styles.dateAndTime}>{formatDateFromIso(row.actionTstampClk, 'DD MMM YY HH:mm')}</div>,
    { id: 'move-history-date-time' },
  ),
  createHeader(
    'Event',
    (row) => <div className={styles.event}>{getMoveHistoryEventTemplate(row).getEventNameDisplay(row)}</div>,
    {
      id: 'move-history-event',
    },
  ),
  createHeader(
    'Details',
    (row) => {
      return <MoveHistoryDetailsSelector historyRecord={row} />;
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
  const useGetMoveHistoryQuery = ({ currentPage, currentPageSize }) => {
    return useGHCGetMoveHistory({ moveCode, currentPage, currentPageSize });
  };

  return (
    <div className={styles.MoveHistoryTable}>
      <TableQueue
        showFilters={false}
        showPagination
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
