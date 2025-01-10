import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.reviewShipmentAddressUpdate,
  tableName: t.shipment_address_updates,
  getEventNameDisplay: () => {
    return 'Shipment Destination Address Request';
  },
  getDetails: ({ changedValues }) => {
    if (changedValues.status === 'APPROVED') {
      return (
        <div>
          <b>Status</b>: Approved
        </div>
      );
    }
    if (changedValues.status === 'REJECTED') {
      return (
        <div>
          <b>Status</b>: Rejected
        </div>
      );
    }
    return null;
  },
};
