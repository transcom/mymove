import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: null,
  eventName: null,
  tableName: null,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: ({ tableName }) => {
    switch (tableName) {
      case t.orders:
        return 'juan Updated order';
      case t.mto_service_items:
        return 'juan Updated service item';
      case t.entitlements:
        return 'juan Updated allowances';
      case t.payment_requests:
        return 'juan Updated payment request';
      case t.mto_shipments:
      case t.mto_agents:
      case t.addresses:
        return 'juan Updated shipment';
      case t.moves:
      default:
        return 'juan Updated move';
    }
  },
  getDetailsPlainText: () => {
    return '-';
  },
};
