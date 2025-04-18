import React from 'react';

import styles from './acknowledgeMovesAndShipments.module.scss';

import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import e from 'constants/MoveHistory/UIDisplay/eventDisplayNames';

export default {
  action: a.UPDATE,
  eventName: o.acknowledgeMovesAndShipments,
  tableName: t.moves,
  getEventNameDisplay: () => e.UPDATED_MOVE,
  getDetails: (historyRecord) => {
    const primeAcknowledgedAt = historyRecord?.changedValues?.prime_acknowledged_at;
    return (
      <>
        <span className={styles.field}>Prime Acknowledged At: </span>
        <span>{primeAcknowledgedAt}</span>
      </>
    );
  },
};
