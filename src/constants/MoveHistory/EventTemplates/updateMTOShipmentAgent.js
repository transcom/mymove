import { formatMoveHistoryAgent } from 'utils/formatters';
import o from 'constants/MoveHistory/UIDisplay/Operations';
import d from 'constants/MoveHistory/UIDisplay/DetailsTypes';
import a from 'constants/MoveHistory/Database/Actions';
import t from 'constants/MoveHistory/Database/Tables';

export default {
  action: a.UPDATE,
  eventName: o.updateMTOShipment,
  tableName: t.mto_agents,
  detailsType: d.LABELED,
  getEventNameDisplay: () => 'Updated shipment',
  getDetailsLabeledDetails: ({ oldValues, changedValues, context }) => {
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
      shipment_type: context[0].shipment_type,
      ...changedValues,
    };

    newChangedValues[agentLabel] = agent;

    return newChangedValues;
  },
};
