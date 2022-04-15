import React from 'react';

import PlainTextDetails from './PlainTextDetails';
import LabeledDetails from './LabeledDetails';

import {
  eventNamesWithLabeledDetails,
  eventNamesWithServiceItemDetails,
  eventNamesWithPlainTextDetails,
  HistoryLogRecordShape,
} from 'constants/historyLogUIDisplayName';

const formatChangedValues = (values) => {
  return values
    ? Object.keys(values).map((key) => (
        <div key={`${key}-${values[key]}`}>
          {key}: {values[key]}
        </div>
      ))
    : '';
};

const MoveHistoryDetailsSelector = ({ historyRecord }) => {
  if (eventNamesWithLabeledDetails[historyRecord.eventName]) {
    return <LabeledDetails changedValues={historyRecord.changedValues} />;
  }

  if (eventNamesWithServiceItemDetails[historyRecord.eventName]) {
    return (
      <div>
        Service Items {historyRecord.eventName}
        <div>old Values {formatChangedValues(historyRecord.oldValues)}</div>
        <div>changed values {formatChangedValues(historyRecord.changedValues)}</div>
      </div>
    );
  }

  if (eventNamesWithPlainTextDetails[historyRecord.eventName]) {
    return <PlainTextDetails historyRecord={historyRecord} />;
  }

  return (
    <div>
      - {historyRecord.eventName}
      <div>old Values {formatChangedValues(historyRecord.oldValues)}</div>
      <div>changed values {formatChangedValues(historyRecord.changedValues)}</div>
    </div>
  );
};

MoveHistoryDetailsSelector.propTypes = {
  historyRecord: HistoryLogRecordShape,
};

MoveHistoryDetailsSelector.defaultProps = {
  historyRecord: {},
};

export default MoveHistoryDetailsSelector;
