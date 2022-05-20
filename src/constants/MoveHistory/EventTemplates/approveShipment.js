import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import t from 'constants/MoveHistory/Database/Tables';
import { shipmentTypes as s } from 'constants/shipments';

export default {
  action: a.UPDATE,
  eventName: o.approveShipment,
  tableName: t.mto_shipments,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => 'Approved shipment',
  getDetailsPlainText: (historyRecord) => {
    return `${s[historyRecord.oldValues?.shipment_type]} shipment`;
  },
};
