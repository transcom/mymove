import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const { changedValues } = historyRecord;
  const { emailIsPreferred, phoneIsPreferred } = changedValues;
  const newChangedValues = { ...changedValues };
  if (emailIsPreferred) {
    newChangedValues.email_is_preferred = 'Yes';
  } else {
    newChangedValues.email_is_preferred = 'No';
  }
  if (phoneIsPreferred) {
    newChangedValues.phone_is_preferred = 'Yes';
  } else {
    newChangedValues.phone_is_preferred = 'No';
  }
  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.updateCustomer,
  tableName: t.service_members,
  getEventNameDisplay: () => 'Updated profile',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
