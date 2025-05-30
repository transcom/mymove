import React from 'react';

import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  let newChangedValues = {
    ...historyRecord.changedValues,
  };

  if (historyRecord.context) {
    newChangedValues = {
      ...newChangedValues,
      ...historyRecord.context[0],
    };
  }

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.reviewShipmentAddressUpdate,
  tableName: t.moves,
  getEventNameDisplay: () => 'Updated move',
  getDetails: (historyRecord) => (
    <>
      <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />
      {historyRecord?.changedValues?.too_task_order_assigned_id !== undefined && (
        <div>Task ordering officer unassigned</div>
      )}
      {historyRecord?.changedValues?.too_destination_assigned_id !== undefined && (
        <div>Destination task ordering officer unassigned</div>
      )}
    </>
  ),
};
