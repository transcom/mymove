import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: '*',
  tableName: t.orders,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Updated orders',
  getDetailsLabeledDetails: (historyRecord) => {
    let newChangedValues;

    if (historyRecord.context) {
      newChangedValues = {
        ...historyRecord.changedValues,
        ...historyRecord.context[0],
      };
    } else {
      newChangedValues = historyRecord.changedValues;
    }

    // merge context with change values for only this event
    return newChangedValues;
  },
};
