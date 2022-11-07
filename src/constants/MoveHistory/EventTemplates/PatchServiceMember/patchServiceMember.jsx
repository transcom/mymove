import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

export default {
  action: a.UPDATE,
  eventName: o.patchServiceMember,
  tableName: t.service_members,
  getEventNameDisplay: () => 'Updated profile',
  getDetails: (historyRecord) => {
    if (historyRecord.changedValues.duty_location_id) {
      const { changedValues } = historyRecord;
      const dutyLocationName = historyRecord.context[0].current_duty_location_name;
      const newChangedValues = {
        ...changedValues,
        current_duty_location_name: dutyLocationName,
      };
      return <LabeledDetails historyRecord={{ ...historyRecord, changedValues: newChangedValues }} />;
    }
    return <LabeledDetails historyRecord={historyRecord} />;
  },
};
