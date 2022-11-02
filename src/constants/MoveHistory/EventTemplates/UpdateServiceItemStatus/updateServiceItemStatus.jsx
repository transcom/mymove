import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { shipmentTypes as s } from 'constants/shipments';

export default {
  action: a.UPDATE,
  eventName: o.updateServiceItemStatus,
  tableName: t.mto_service_items,
  getEventNameDisplay: ({ changedValues }) => {
    return changedValues?.status === 'APPROVED' ? <> Approved service item </> : <> Rejected service item </>;
  },
  getDetails: ({ context }) => (
    <>
      {s[context[0]?.shipment_type]} shipment #{context[0]?.shipment_id_abbr.toUpperCase()}, {context[0]?.name}
    </>
  ),
};
