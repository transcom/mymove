import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipmentStatus,
  tableName: t.mto_shipments,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsLabeledDetails: ({ changedValues, oldValues }) => {
    return {
      shipment_type: oldValues.shipment_type,
      ...changedValues,
    };
  },
};
