import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import { shipmentTypes as s } from 'constants/shipments';

export default {
  action: '*',
  eventName: o.approveShipmentDiversion,
  tableName: '*',
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => 'Approved diversion',
  getDetailsPlainText: (historyRecord) => {
    return `${s[historyRecord.oldValues?.shipment_type]} shipment`;
  },
};
