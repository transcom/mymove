import React from 'react';

import styles from './MultiMovesMoveInfoList.module.scss';

import descriptionListStyles from 'styles/descriptionList.module.scss';
import { formatDateForDatePicker } from 'shared/dates';

const MultiMovesMoveInfoList = ({ move }) => {
  const { orders } = move;

  // function that determines label based on order type
  const getReportByLabel = (ordersType) => {
    if (ordersType === 'SEPARATION') {
      return 'Separation Date';
    }
    if (ordersType === 'RETIREMENT') {
      return 'Retirement Date';
    }
    return 'Report by Date';
  };

  // function that determines label based on order type
  const getOrdersTypeLabel = (ordersType) => {
    if (ordersType === 'SEPARATION') {
      return 'Separation';
    }
    if (ordersType === 'RETIREMENT') {
      return 'Retirement';
    }
    return 'Permanent Change of Station';
  };

  // destination duty location label will differ based on order type
  const getDestinationDutyLocationLabel = (ordersType) => {
    if (ordersType === 'SEPARATION') {
      return 'HOR or PLEAD';
    }
    if (ordersType === 'RETIREMENT') {
      return 'HOR, HOS, or PLEAD';
    }
    return 'Destination Duty Location';
  };

  const toCamelCase = (str) => {
    return str.replace(/_([a-z])/g, (match, letter) => letter.toUpperCase());
  };

  const formatAddress = (address) => {
    const camelCaseAddress = Object.keys(address).reduce((acc, key) => {
      acc[toCamelCase(key)] = address[key];
      return acc;
    }, {});

    const { streetAddress1, streetAddress2, streetAddress3, city, state, postalCode, id } = camelCaseAddress;

    // Check for empty UUID
    const isIdEmpty = id === '00000000-0000-0000-0000-000000000000';

    // Check for null values and empty UUID
    if (isIdEmpty) {
      return '-';
    }

    return (
      <>
        {streetAddress1 && <>{streetAddress1},&nbsp;</>}
        {streetAddress2 && <>{streetAddress2},&nbsp;</>}
        {streetAddress3 && <>{streetAddress3},&nbsp;</>}
        {city ? `${city}, ${state} ${postalCode}` : postalCode}
      </>
    );
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
            <dd>{formatDateForDatePicker(orders.issue_date) || '-'}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Orders Type</dt>
            <dd>{getOrdersTypeLabel(orders.orders_type) || '-'}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>{getReportByLabel(orders.ordersType)}</dt>
            <dd>{formatDateForDatePicker(orders.report_by_date) || '-'}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>Current Duty Location</dt>
            <dd>{formatAddress(orders.OriginDutyLocation.Address) || '-'}</dd>
          </div>

          <div className={descriptionListStyles.row}>
            <dt>{getDestinationDutyLocationLabel(orders.ordersType)}</dt>
            <dd>{formatAddress(orders.NewDutyLocation.Address) || '-'}</dd>
          </div>
        </dl>
      </div>
    </div>
  );
};

export default MultiMovesMoveInfoList;
