import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipment,
  tableName: t.mto_shipments,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsLabeledDetails: (historyRecord) => {
    return {
      shipment_type: historyRecord.oldValues.shipment_type,
      // TODO: [ MB-12182 ] This will include a shipment ID label in the future
      ...historyRecord.changedValues,
    };
  },
};
