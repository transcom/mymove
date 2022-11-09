import React from 'react';

import { formatMoveHistoryAgent } from 'utils/formatters';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const { oldValues, changedValues } = historyRecord;
  let newChangedValues = {
    email: oldValues.email,
    first_name: oldValues.first_name,
    last_name: oldValues.last_name,
    phone: oldValues.phone,
    ...changedValues,
  };

  const agent = formatMoveHistoryAgent(newChangedValues);

  const agentType = changedValues.agent_type ?? oldValues.agent_type;

  let agentLabel = '';
  if (agentType === 'RECEIVING_AGENT') {
    agentLabel = 'receiving_agent';
  } else if (agentType === 'RELEASING_AGENT') {
    agentLabel = 'releasing_agent';
  }

  newChangedValues = {
    ...changedValues,
    ...getMtoShipmentLabel(historyRecord),
  };

  newChangedValues[agentLabel] = agent;

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipment,
  tableName: t.mto_agents,
  getEventNameDisplay: () => 'Updated shipment',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
