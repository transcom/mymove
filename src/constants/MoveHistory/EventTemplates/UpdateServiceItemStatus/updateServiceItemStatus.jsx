import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { shipmentTypes as s } from 'constants/shipments';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

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
    const { context, changedValues } = historyRecord;
    const formattedContext = {
      shipment_type: s[context[0]?.shipment_type],
      shipment_id_display: context[0]?.shipment_id_abbr.toUpperCase(),
      service_item_name: context[0]?.name,
    };

    return changedValues?.status === 'REJECTED' ? (
      <LabeledDetails historyRecord={formatChangedValues(historyRecord, formattedContext)} />
    ) : (
      <>
        {formattedContext.shipment_type} shipment #{formattedContext.shipment_id_display},{' '}
        {formattedContext.service_item_name}
      </>
    );
  },
};
