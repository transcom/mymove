import detailsTypes from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import moveHistoryOperations from 'constants/MoveHistory/UIDisplay/Operations';
import dbActions from 'constants/MoveHistory/Database/Actions';
import dbTables from 'constants/MoveHistory/Database/Tables';

export default {
  action: dbActions.UPDATE,
  eventName: moveHistoryOperations.updateMTOShipment,
  tableName: dbTables.mto_shipments,
  detailsType: detailsTypes.LABELED,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsLabeledDetails: (historyRecord) => {
    return {
      shipment_type: historyRecord.oldValues.shipment_type,
      ...historyRecord.changedValues,
    };
  },
};
