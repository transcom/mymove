import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './PaymentDetails.module.scss';

import { HistoryLogRecordShape } from 'constants/historyLogUIDisplayName';

const iconToDisplay = (statusToFilter) => {
  if (statusToFilter === 'Approved') {
    return <FontAwesomeIcon icon="check" className={styles.successCheck} />;
  }
  return <FontAwesomeIcon icon="times" className={styles.rejectTimes} />;
};

const filterContextStatus = (context, statusToFilter) => {
  const contextToDisplay = [];
  let sum = 0;
  context.forEach((value) => {
    if (value.status === statusToFilter.toUpperCase()) {
      const price = parseFloat(value.price);
      sum += price;
      contextToDisplay.push(
        <div className={styles.serviceItemRow} key={`${value.name}`}>
          <div>{value.name}</div>
          <div>{price.toFixed(2)}</div>
        </div>,
      );
    }
  });
  return (
    <div>
      <div className={styles.statusRow}>
        <b>{statusToFilter} service items total: </b>
        <div>
          {iconToDisplay(statusToFilter)} &nbsp;
          <b>${sum.toFixed(2)}</b>
        </div>
      </div>
      <div>{contextToDisplay}</div>
    </div>
  );
};

const PaymentDetails = ({ historyRecord }) => {
  return (
    <div className={styles.PaymentDetails}>
      {filterContextStatus(historyRecord.context, 'Approved')}
      {filterContextStatus(historyRecord.context, 'Rejected')}
    </div>
  );
};

PaymentDetails.propTypes = {
  historyRecord: HistoryLogRecordShape,
};

PaymentDetails.defaultProps = {
  historyRecord: {},
};

export default PaymentDetails;
