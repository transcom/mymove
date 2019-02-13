import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';

const getServiceAgentsForShipmentLabel = 'ServiceAgents.getServiceAgentsForShipment';
const updateServiceAgentForShipmentLabel = 'ServiceAgents.updateServiceAgentForShipment';

export function getServiceAgentsForShipment(shipmentId) {
  const label = getServiceAgentsForShipmentLabel;
  const swaggerTag = 'service_agents.indexServiceAgents';
  return swaggerRequest(getPublicClient, swaggerTag, { shipmentId }, { label });
}

export function updateServiceAgentForShipment(shipmentId, serviceAgentId, serviceAgent) {
  const label = updateServiceAgentForShipmentLabel;
  const swaggerTag = 'service_agents.patchServiceAgent';
  return swaggerRequest(getPublicClient, swaggerTag, { shipmentId, serviceAgentId, serviceAgent }, { label });
}

export function updateServiceAgentsForShipment(shipmentId, serviceAgents) {
  return async function(dispatch) {
    Object.values(serviceAgents).map(serviceAgent =>
      dispatch(updateServiceAgentForShipment(shipmentId, serviceAgent.id, serviceAgent)),
    );
  };
}

export function selectServiceAgentsForShipment(state, shipmentId) {
  if (!shipmentId) {
    return [];
  }
  const serviceAgents = Object.values(state.entities.serviceAgents);
  return serviceAgents.filter(serviceAgent => serviceAgent.shipment_id === shipmentId);
}
