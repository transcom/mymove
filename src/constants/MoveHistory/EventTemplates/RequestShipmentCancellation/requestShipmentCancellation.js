import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import { shipmentTypes } from 'constants/shipments';

export default {
  action: '*',
  eventName: o.requestShipmentCancellation,
  tableName: '*',
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsPlainText: (historyRecord) => {
    return `Requested cancellation for ${shipmentTypes[historyRecord.oldValues?.shipment_type]} shipment`;
  },
};
