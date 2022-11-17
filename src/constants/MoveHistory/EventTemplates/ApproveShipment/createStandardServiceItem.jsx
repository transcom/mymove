import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { shipmentTypes as SHIPMENT_TYPE } from 'constants/shipments';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';

export default {
  action: a.INSERT,
  eventName: o.approveShipment,
  tableName: t.mto_service_items,
  getEventNameDisplay: () => 'Approved service item',
  getDetails: (historyRecord) => {
    const shipmentLabel = getMtoShipmentLabel(historyRecord);
    return (
      <>
        {SHIPMENT_TYPE[shipmentLabel.shipment_type]} shipment #{shipmentLabel.shipment_id_display}
        {', '}
        {shipmentLabel.service_item_name}
      </>
    );
  },
};
