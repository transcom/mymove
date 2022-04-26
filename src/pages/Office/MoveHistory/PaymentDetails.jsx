import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './PaymentDetails.module.scss';

import { HistoryLogContextShape } from 'constants/historyLogUIDisplayName';

const APPROVED_STRING = 'Approved';
const REJECTED_STRING = 'Rejected';

const iconToDisplay = (statusToFilter) => {
  if (statusToFilter === APPROVED_STRING) {
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

const PaymentDetails = ({ context }) => {
  return (
    <div className={styles.PaymentDetails}>
      {filterContextStatus(context, APPROVED_STRING)}
      {filterContextStatus(context, REJECTED_STRING)}
    </div>
  );
};

PaymentDetails.propTypes = {
  context: HistoryLogContextShape,
};

PaymentDetails.defaultProps = {
  context: {},
};

export default PaymentDetails;
