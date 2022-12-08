import React from 'react';

import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const { changedValues } = historyRecord;
  const { name, phone, email } = changedValues;
  const newChangedValues = { ...changedValues };
  if (name) {
    newChangedValues.backup_contact_name = changedValues.name;
    delete newChangedValues.name;
  }
  if (email) {
    newChangedValues.backup_contact_email = changedValues.email;
    delete newChangedValues.email;
  }
  if (phone) {
    newChangedValues.backup_contact_phone = changedValues.phone;
    delete newChangedValues.phone;
  }
  return { ...historyRecord, changedValues: newChangedValues };
};
export default {
  action: a.UPDATE,
  eventName: '*', // Needs wild card to handle both updateServiceMemberBackupContact and updateCustomer event
  tableName: t.backup_contacts,
  getEventNameDisplay: () => 'Updated profile',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
