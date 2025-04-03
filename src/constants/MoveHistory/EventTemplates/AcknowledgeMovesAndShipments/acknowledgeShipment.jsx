import React from 'react';

import styles from './acknowledgeMovesAndShipments.module.scss';

import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import e from 'constants/MoveHistory/UIDisplay/eventDisplayNames';
import { shipmentTypes as s } from 'constants/shipments';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';

export default {
  action: a.UPDATE,
  eventName: o.acknowledgeMovesAndShipments,
  tableName: t.mto_shipments,
  getEventNameDisplay: () => e.UPDATED_SHIPMENT,
  getDetails: (historyRecord) => {
    const primeAcknowledgedAt = historyRecord?.changedValues?.prime_acknowledged_at;
    const formattedContext = getMtoShipmentLabel(historyRecord);
    return (
      <>
        <span className={styles.field}>Prime Acknowledged At: </span>
        <span>{primeAcknowledgedAt}</span>
        <div>
          {s[formattedContext.shipment_type]} shipment #{formattedContext.shipment_locator}
        </div>
      </>
    );
  },
};
