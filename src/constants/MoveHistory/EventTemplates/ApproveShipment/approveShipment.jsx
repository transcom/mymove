import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { shipmentTypes as s } from 'constants/shipments';

export default {
  action: a.UPDATE,
  eventName: '*', // both approveShipment and approveShipments events can render this template
  tableName: t.mto_shipments,
  getEventNameDisplay: () => 'Approved shipment',
  getDetails: ({ context }) => (
    <>
      {s[context[0]?.shipment_type]} shipment #{context[0]?.shipment_locator.toUpperCase()}
    </>
  ),
};
