import React from 'react';

import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

export default {
  action: a.INSERT,
  eventName: o.createOrders,
  tableName: t.orders,
  getEventNameDisplay: () => 'Created orders',
  getDetails: (historyRecord) => {
    const getDetailsLabeledDetails = ({ context, changedValues }) => {
      const newChangedValues = {
        new_duty_location_name: context[0]?.new_duty_location_name,
        origin_duty_location_name: context[0]?.origin_duty_location_name,
        ...changedValues,
      };

      return newChangedValues;
    };

    return <LabeledDetails historyRecord={historyRecord} getDetailsLabeledDetails={getDetailsLabeledDetails} />;
  },
};
