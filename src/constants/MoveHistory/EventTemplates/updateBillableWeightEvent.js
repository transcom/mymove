import detailsTypes from 'constants/MoveHistory/DetailsColumn/Types';
import dbTables from 'constants/MoveHistory/Database/Tables';

export default {
  action: null,
  eventName: null,
  tableName: null,
  detailsType: detailsTypes.PLAIN_TEXT,
  getEventNameDisplay: ({ tableName }) => {
    switch (tableName) {
      case dbTables.orders:
        return 'Updated order';
      case dbTables.mto_service_items:
        return 'Updated service item';
      case dbTables.entitlements:
        return 'Updated allowances';
      case dbTables.payment_requests:
        return 'Updated payment request';
      case dbTables.mto_shipments:
      case dbTables.mto_agents:
      case dbTables.addresses:
        return 'Updated shipment';
      case dbTables.moves:
      default:
        return 'Updated move';
    }
  },
  getDetailsPlainText: () => {
    return '-';
  },
};
