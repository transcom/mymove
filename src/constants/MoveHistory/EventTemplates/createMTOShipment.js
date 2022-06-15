import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.INSERT,
  eventName: o.createMTOShipment,
  tableName: t.mto_shipments,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Created shipment',
};
