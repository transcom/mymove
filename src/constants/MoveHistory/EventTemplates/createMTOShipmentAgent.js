import { formatMoveHistoryAgent } from 'utils/formatters';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import { getMtoShipmentLabel } from 'utils/formatMtoShipment';

export default {
  action: a.INSERT,
  eventName: '*',
  tableName: t.mto_agents,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsLabeledDetails: ({ changedValues, oldValues, context }) => {
    const agent = formatMoveHistoryAgent(changedValues);

    const agentType = changedValues.agent_type ?? oldValues.agent_type;

    let agentLabel = '';
    if (agentType === 'RECEIVING_AGENT') {
      agentLabel = 'receiving_agent';
    } else if (agentType === 'RELEASING_AGENT') {
      agentLabel = 'releasing_agent';
    }

    const newChangedValues = {
      ...changedValues,
      ...getMtoShipmentLabel({ context }),
    };

    newChangedValues[agentLabel] = agent;

    return newChangedValues;
  },
};
