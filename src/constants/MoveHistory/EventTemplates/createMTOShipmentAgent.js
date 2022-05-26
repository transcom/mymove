import { formatMoveHistoryAgent } from 'utils/formatters';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import t from 'constants/MoveHistory/Database/Tables';
import a from 'constants/MoveHistory/Database/Actions';
import o from 'constants/MoveHistory/UIDisplay/Operations';

export default {
  action: a.INSERT,
  eventName: o.updateMTOShipment,
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
      shipment_type: context[0]?.shipment_type,
      ...changedValues,
    };

    newChangedValues[agentLabel] = agent;

    return newChangedValues;
  },
};
