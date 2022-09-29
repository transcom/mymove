import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updateMTOStatusServiceCounselingCompleted,
  tableName: t.moves,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => `Updated shipment`,
  getDetailsPlainText: () => {
    return `Counseling Completed`;
  },
};
