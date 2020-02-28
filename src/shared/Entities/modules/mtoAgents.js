import { getGHCClient } from 'shared/Swagger/api';
import { swaggerRequest } from 'shared/Swagger/request';
import { get } from 'lodash';

const getMTOAgentListOperation = 'mtoAgent.fetchMTOAgentList';
const mtoAgentsSchemaKey = 'mtoAgent';

export function getMTOAgentList(
  moveTaskOrderID,
  shipmentID,
  label = getMTOAgentListOperation,
  schemaKey = mtoAgentsSchemaKey,
) {
  return swaggerRequest(getGHCClient, getMTOAgentListOperation, { moveTaskOrderID, shipmentID }, { label, schemaKey });
}

export function selectMTOAgents(state) {
  const mtoAgents = get(state, 'entities.mtoAgents') || {};
  return Object.values(mtoAgents);
}
