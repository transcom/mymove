import React from 'react';

import styles from './MultiMovesMoveInfoList.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatAddress } from 'utils/shipmentDisplay';

const MultiMovesMoveInfoList = ({ move }) => {
  const { orders } = move;

  return (
    <div className={styles.moveInfoContainer} data-testid="move--info-container">
      <div className={styles.moveInfoSection}>
        <dl className={descriptionListStyles.descriptionList}>
          <div className={descriptionListStyles.row}>
            <dt>Move Status</dt>
            <dd>{move.status || '-'}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Orders Issue Date</dt>
            <dd>{orders.date_issued || '-'}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Orders Type</dt>
            <dd>{orders.ordersType || '-'}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Report by</dt>
            <dd>{orders.reportByDate || '-'}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Current Duty Location</dt>
            <dd>{formatAddress(orders.originDutyLocation.address) || '-'}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Destination Duty Location</dt>
            <dd>{formatAddress(orders.destinationDutyLocation.address) || '-'}</dd>
          </div>
        </dl>
      </div>
    </div>
  );
};

export default MultiMovesMoveInfoList;
