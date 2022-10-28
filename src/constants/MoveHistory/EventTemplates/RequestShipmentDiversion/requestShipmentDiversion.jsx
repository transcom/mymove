import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import { shipmentTypes as s } from 'constants/shipments';

export default {
  action: a.UPDATE,
  eventName: o.requestShipmentDiversion,
  tableName: t.mto_shipments,
  getEventNameDisplay: () => 'Requested diversion',
  getDetails: ({ context }) => (
    <>
      Requested diversion for {s[context[0]?.shipment_type]} shipment #{context[0].shipment_id_abbr.toUpperCase()}
    </>
  ),
};
