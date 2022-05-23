import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { shipmentTypes } from 'constants/shipments';

export default {
  action: a.INSERT,
  eventName: o.approveShipment,
  tableName: t.mto_service_items,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => 'Approved service item',
  getDetailsPlainText: (historyRecord) => {
    return `${shipmentTypes[historyRecord.context[0]?.shipment_type]} shipment, ${historyRecord.context[0]?.name}`;
  },
};
