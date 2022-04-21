import React from 'react';

import LabeledDetails from './LabeledDetails';

import { HistoryLogRecordShape } from 'constants/historyLogUIDisplayName';
import getMoveHistoryEventTemplate, { detailsTypes } from 'constants/moveHistoryEventTemplate';

const MoveHistoryDetailsSelector = ({ historyRecord }) => {
  const eventTemplate = getMoveHistoryEventTemplate(historyRecord);
  const changedValuesMerged = Object.assign(historyRecord.changedValues, historyRecord.context);
  switch (eventTemplate.detailsType) {
    case detailsTypes.LABELED:
      return <LabeledDetails changedValues={changedValuesMerged} />;
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
