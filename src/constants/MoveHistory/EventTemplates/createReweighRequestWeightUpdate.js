import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { shipmentTypes } from 'constants/shipments';

export default {
  action: a.INSERT,
  eventName: o.updateMTOShipment,
  tableName: t.reweighs,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => `Updated shipment`,
  getDetailsPlainText: ({ context }) => {
    const shipmentType = context[0].shipment_type;
    return `${shipmentTypes[shipmentType]} shipment, reweigh requested`;
  },
};
