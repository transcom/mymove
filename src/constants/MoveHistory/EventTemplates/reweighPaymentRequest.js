import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updateReweigh,
  tableName: t.reweighs,
  detailsType: d.STATUS,
  getEventNameDisplay: ({ context }) => `Updated payment request ${context[0]?.payment_request_number}`,
  getStatusDetails: () => {
    return 'Recalculated payment request';
  },
};
