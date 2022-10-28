import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import { shipmentTypes as s } from 'constants/shipments';

export default {
  action: '*',
  eventName: o.approveShipmentDiversion,
  tableName: '*',
  getEventNameDisplay: () => 'Approved diversion',
  getDetails: ({ context }) => (
    <>
      {s[context[0].shipment_type]} shipment #{context[0].shipment_id_abbr.toUpperCase()}
    </>
  ),
};
