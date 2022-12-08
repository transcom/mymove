import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { shipmentTypes as s } from 'constants/shipments';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';

const formatChangedValues = (historyRecord, formattedContext) => {
  const { changedValues } = historyRecord;

  const newChangedValues = {
    ...formattedContext,
    rejection_reason: changedValues.rejection_reason,
    rejected_at: changedValues.rejected_at,
  };

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.updateServiceItemStatus,
  tableName: t.mto_service_items,
  getEventNameDisplay: ({ changedValues }) => {
    return changedValues?.status === 'APPROVED' ? <> Approved service item </> : <> Rejected service item </>;
  },
  getDetails: (historyRecord) => {
    const { changedValues } = historyRecord;
    const formattedContext = getMtoShipmentLabel(historyRecord);

    return changedValues?.status === 'REJECTED' ? (
      <LabeledDetails historyRecord={formatChangedValues(historyRecord, formattedContext)} />
    ) : (
      <>
        {s[formattedContext.shipment_type]} shipment #{formattedContext.shipment_id_display},{' '}
        {formattedContext.service_item_name}
      </>
    );
  },
};
