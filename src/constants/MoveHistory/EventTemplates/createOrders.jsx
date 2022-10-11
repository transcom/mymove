import React from 'react';

import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

export default {
  action: a.INSERT,
  eventName: o.createOrders,
  tableName: t.orders,
  detailsType: d.PLAIN_TEXT,
  getEventNameDisplay: () => 'Created orders',
  getDetails: (historyRecord) => {
    const getDetailsLabeledDetails = (history) => {
      const newChangedValues = {
        new_duty_location_name: history.context[0]?.new_duty_location_name,
        origin_duty_location_name: history.context[0]?.origin_duty_location_name,
        ...history.changedValues,
      };

      return newChangedValues;
    };

    return <LabeledDetails historyRecord={historyRecord} getDetailsLabeledDetails={getDetailsLabeledDetails} />;
  },
};
