import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import { shipmentTypes as s } from 'constants/shipments';

export default {
  action: a.UPDATE,
  eventName: o.deleteShipment,
  tableName: t.mto_shipments,
  getEventNameDisplay: () => 'Deleted shipment',
  getDetails: ({ context }) => (
    <>
      {s[context[0].shipment_type]} shipment #{context[0].shipment_id_abbr.toUpperCase()} deleted
    </>
  ),
};
