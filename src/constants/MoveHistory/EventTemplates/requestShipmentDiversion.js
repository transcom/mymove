import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import { shipmentTypes } from 'constants/shipments';

export default {
  action: '*',
  eventName: o.requestShipmentDiversion,
  tableName: '*',
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => 'Requested diversion',
  getDetailsPlainText: (historyRecord) => {
    return `Requested diversion for ${shipmentTypes[historyRecord.oldValues?.shipment_type]} shipment`;
  },
};
