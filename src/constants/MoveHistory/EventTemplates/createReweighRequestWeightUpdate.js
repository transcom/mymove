import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.INSERT,
  eventName: o.updateMTOShipment,
  tableName: t.reweighs,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => `Updated shipment`,
  getDetailsPlainText: ({ context }) => {
    return `${context[0].shipment_type} shipment, reweigh requested`;
  },
};
