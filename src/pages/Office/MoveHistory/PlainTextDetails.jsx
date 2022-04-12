import React from 'react';

import { detailsPlainTextToDisplay, HistoryLogRecordShape } from 'constants/historyLogUIDisplayName';

const PlainTextDetails = ({ historyRecord }) => {
  return <div>{detailsPlainTextToDisplay(historyRecord)}</div>;
};

PlainTextDetails.propTypes = {
  historyRecord: HistoryLogRecordShape,
};

PlainTextDetails.defaultProps = {
  historyRecord: {},
};

export default PlainTextDetails;
