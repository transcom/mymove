import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { formatCents, formatWeight } from 'utils/formatters';

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipment,
  tableName: t.mto_service_items,
  getEventNameDisplay: ({ changedValues }) => {
    if (changedValues.pricing_estimate) {
      return 'Service item estimated price updated';
    }
    if (changedValues.estimated_weight) {
      return 'Service item estimated weight updated';
    }
    return 'Service item updated';
  },
  getDetails: ({ changedValues, context }) => {
    if (changedValues.pricing_estimate) {
      return (
        <div>
          <b>Service item</b>: {context[0].name}
          <br />
          <b>Estimated Price</b>: ${formatCents(changedValues.pricing_estimate)}
        </div>
      );
    }
    if (changedValues.estimated_weight) {
      return (
        <div>
          <b>Service item</b>: {context[0].name}
          <br />
          <b>Estimated weight</b>: {formatWeight(changedValues.estimated_weight)}
        </div>
      );
    }
    return null;
  },
};
