import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import { shipmentTypes as s } from 'constants/shipments';

export default {
  action: a.INSERT,
  eventName: o.requestShipmentReweigh,
  tableName: t.reweighs,
  getEventNameDisplay: () => 'Updated shipment',
  getDetails: ({ context }) => (
    <>
      {s[context[0].shipment_type]} shipment #{context[0].shipment_id_abbr.toUpperCase()}, reweigh requested
    </>
  ),
};
