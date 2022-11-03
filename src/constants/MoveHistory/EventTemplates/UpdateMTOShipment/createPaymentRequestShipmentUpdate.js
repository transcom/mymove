import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.INSERT,
  eventName: o.updateMTOShipment,
  tableName: t.payment_requests,
  detailsType: d.STATUS,
  getEventNameDisplay: ({ changedValues }) => `Created payment request ${changedValues?.payment_request_number}`,
  getStatusDetails: () => {
    return 'Pending';
  },
};
