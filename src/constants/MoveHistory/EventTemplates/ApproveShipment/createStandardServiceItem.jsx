import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { shipmentTypes as s } from 'constants/shipments';

export default {
  action: a.INSERT,
  eventName: o.approveShipment,
  tableName: t.mto_service_items,
  getEventNameDisplay: () => 'Approved service item',
  getDetails: ({ context }) => (
    <>
      {s[context[0]?.shipment_type]} shipment #{context[0].shipment_id_abbr.toUpperCase()}, {context[0]?.name}
    </>
  ),
};
