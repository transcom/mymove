import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.createMTOServiceItem,
  tableName: t.moves,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Updated move',
  getDetailsLabeledDetails: (historyRecord) => {
    return {
      status: historyRecord.oldValues.status,
      ...historyRecord.changedValues,
    };
  },
};
