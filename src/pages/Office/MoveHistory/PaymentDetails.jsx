import React from 'react';

import { HistoryLogRecordShape } from 'constants/historyLogUIDisplayName';

const PaymentDetails = ({ historyRecord }) => {
  return <div>{historyRecord}</div>;
};

PaymentDetails.propTypes = {
  historyRecord: HistoryLogRecordShape,
};

PaymentDetails.defaultProps = {
  historyRecord: {},
};

export default PaymentDetails;
