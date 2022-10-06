import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import { shipmentTypes as s } from 'constants/shipments';

export default {
  action: a.UPDATE,
  eventName: o.deleteShipment,
  tableName: t.mto_shipments,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => 'Deleted Shipment',
  getDetailsPlainText: (historyRecord) => {
    // TODO: [ MB-12182 ] This will include a shipment ID label in the future
    return `${s[historyRecord.oldValues?.shipment_type]} shipment deleted`;
  },
};
