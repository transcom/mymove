import React from 'react';

import { eventNamePlainTextToDisplay, HistoryLogRecordShape } from 'constants/historyLogUIDisplayName';

const PlainTextDetails = ({ historyRecord }) => {
  let textToDisplay = '';
  if (eventNamePlainTextToDisplay[historyRecord.eventName]) {
    textToDisplay = eventNamePlainTextToDisplay[historyRecord.eventName](historyRecord);
  }
  return <div>{textToDisplay}</div>;
};

PlainTextDetails.propTypes = {
  historyRecord: HistoryLogRecordShape,
};

PlainTextDetails.defaultProps = {
  historyRecord: {},
};

export default PlainTextDetails;
