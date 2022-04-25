import React from 'react';

import { HistoryLogRecordShape } from 'constants/historyLogUIDisplayName';

const filterContextStatus = (context, statusToFilter) => {
  const contextToDisplay = [];
  let sum = 0;
  context.forEach((value) => {
    if (value.status === statusToFilter.toUpperCase()) {
      sum += parseFloat(value.price);
      contextToDisplay.push(
        <div key={`${value.name}`}>
          {value.name} <br />
          {value.price}
        </div>,
      );
    }
  });
  return (
    <div>
      {statusToFilter} service items total: {sum}
      {contextToDisplay}
    </div>
  );
};

const PaymentDetails = ({ historyRecord }) => {
  return (
    <div>
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
