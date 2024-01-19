import React from 'react';

import styles from './MultiMovesMoveInfoList.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatAddress } from 'utils/shipmentDisplay';

const MultiMovesMoveInfoList = ({ move }) => {
  const { orders } = move;

  const getReportByLabel = (ordersType) => {
    if (ordersType === 'SEPARATION') {
      return 'Separation Date';
    }
    if (ordersType === 'RETIREMENT') {
      return 'Retirement Date';
    }
    return 'Report by Date';
  };

  const getDestinationDutyLocationLabel = (ordersType) => {
    if (ordersType === 'SEPARATION') {
      return 'HOR or PLEAD';
    }
    if (ordersType === 'RETIREMENT') {
      return 'HOR, HOS, or PLEAD';
    }
    return 'Destination Duty Location';
  };

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
            <dt>{getReportByLabel(orders.ordersType)}</dt>
            <dd>{orders.reportByDate || '-'}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Current Duty Location</dt>
            <dd>{formatAddress(orders.originDutyLocation.address) || '-'}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>{getDestinationDutyLocationLabel(orders.ordersType)}</dt>
            <dd>{formatAddress(orders.destinationDutyLocation.address) || '-'}</dd>
          </div>
        </dl>
      </div>
    </div>
  );
};

export default MultiMovesMoveInfoList;
