import React from 'react';

import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import { formatMoveHistoryAgent } from 'utils/formatters';
import LabeledDetails from 'pages/Office/MoveHistory/LabeledDetails';

const formatChangedValues = (historyRecord) => {
  const { changedValues, oldValues, context } = historyRecord;
  const agent = formatMoveHistoryAgent(changedValues);

  const agentType = changedValues.agent_type ?? oldValues.agent_type;

  let agentLabel = '';
  if (agentType === 'RECEIVING_AGENT') {
    agentLabel = 'receiving_agent';
  } else if (agentType === 'RELEASING_AGENT') {
    agentLabel = 'releasing_agent';
  }

  const newChangedValues = {
    shipment_type: context[0]?.shipment_type,
    shipment_id_display: context[0]?.shipment_id_abbr.toUpperCase(),
    ...changedValues,
  };

  newChangedValues[agentLabel] = agent;

  return { ...historyRecord, changedValues: newChangedValues };
};

export default {
  action: a.INSERT,
  eventName: o.updateMTOShipment,
  tableName: t.mto_agents,
  getEventNameDisplay: () => 'Updated shipment',
  getDetails: (historyRecord) => <LabeledDetails historyRecord={formatChangedValues(historyRecord)} />,
};
