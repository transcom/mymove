import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { shipmentTypes } from 'constants/shipments';

export default {
  action: a.INSERT,
  eventName: o.updateMTOShipment,
  tableName: t.reweighs,
  getEventNameDisplay: () => `Updated shipment`,
  getDetails: ({ context }) => {
    const shipmentType = context[0].shipment_type;
    const shipmentIdDisplay = context[0].shipment_id_abbr.toUpperCase();
    return (
      <>
        {shipmentTypes[shipmentType]} shipment #{shipmentIdDisplay}, reweigh requested
      </>
    );
  },
};
