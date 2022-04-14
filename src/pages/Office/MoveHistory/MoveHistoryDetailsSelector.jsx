import React from 'react';

import LabeledDetails from './LabeledDetails';

import { HistoryLogRecordShape } from 'constants/historyLogUIDisplayName';
import getMoveHistoryEventTemplate, { detailsTypes } from 'constants/moveHistoryEventTemplate';

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
  const eventTemplate = getMoveHistoryEventTemplate(historyRecord);
  switch (eventTemplate.detailsType) {
    case detailsTypes.LABELED:
      return <LabeledDetails changedValues={historyRecord.changedValues} />;
    case detailsTypes.LABELED_SERVICE_ITEM:
      return (
        <div>
          Service Items {historyRecord.eventName}
          <div>old Values {formatChangedValues(historyRecord.oldValues)}</div>
          <div>changed values {formatChangedValues(historyRecord.changedValues)}</div>
        </div>
      );
    case detailsTypes.PLAIN_TEXT:
    default:
      return <div>{eventTemplate.getDetailsPlainText(historyRecord)}</div>;
  }
};

MoveHistoryDetailsSelector.propTypes = {
  historyRecord: HistoryLogRecordShape,
};

MoveHistoryDetailsSelector.defaultProps = {
  historyRecord: {},
};

export default MoveHistoryDetailsSelector;
