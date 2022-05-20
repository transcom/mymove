import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import { shipmentTypes } from 'constants/shipments';

export default {
  action: a.INSERT,
  eventName: o.requestShipmentReweigh,
  tableName: t.reweighs,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsPlainText: (historyRecord) => {
    return `${shipmentTypes[historyRecord.context[0]?.shipment_type]} shipment, reweigh requested`;
  },
};
