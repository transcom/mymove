import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import { shipmentTypes as s } from 'constants/shipments';

export default {
  action: a.UPDATE,
  eventName: o.approveShipment,
  tableName: t.mto_shipments,
  getEventNameDisplay: () => 'Approved shipment',
  getDetails: ({ context }) => (
    <>
      {s[context[0]?.shipment_type]} shipment #{context[0]?.shipment_id_abbr.toUpperCase()}
    </>
  ),
};
