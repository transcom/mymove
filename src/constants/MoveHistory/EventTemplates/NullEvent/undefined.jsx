import React from 'react';

import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: null,
  eventName: null,
  tableName: null,
  getEventNameDisplay: ({ tableName }) => {
    switch (tableName) {
      case t.orders:
        return <> Updated order </>;
      case t.mto_service_items:
        return <> Updated service item </>;
      case t.entitlements:
        return <> Updated allowances </>;
      case t.payment_requests:
        return <> Updated payment request </>;
      case t.mto_shipments:
      case t.mto_agents:
      case t.addresses:
        return <> Updated shipment </>;
      case t.moves:
      default:
        return <> Updated move </>;
    }
  },
  getDetails: () => <> - </>,
};
