import React from 'react';

import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

export default {
  action: a.INSERT,
  eventName: o.createOrders,
  tableName: t.entitlements,
  getEventNameDisplay: () => 'Created allowances',
  getDetails: (historyRecord) => {
    const getDetailsLabeledDetails = ({ changedValues }) => {
      const newChangedValues = {
        ...changedValues,
        dependents_authorized: changedValues.dependents_authorized === true ? 'Yes' : 'No',
      };

      return newChangedValues;
    };

    return <LabeledDetails historyRecord={historyRecord} getDetailsLabeledDetails={getDetailsLabeledDetails} />;
  },
};
