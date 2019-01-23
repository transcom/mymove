import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';

const getServiceAgentsForShipmentLabel = 'ServiceAgents.getServiceAgentsForShipment';

export function getServiceAgentsForShipment(shipmentId) {
  const label = getServiceAgentsForShipmentLabel;
  const swaggerTag = 'service_agents.indexServiceAgents';
  return swaggerRequest(getPublicClient, swaggerTag, { shipmentId }, { label });
}

export function selectServiceAgentsForShipment(state, shipmentId) {
  if (!shipmentId) {
    return [];
  }
  const serviceAgents = Object.values(state.entities.serviceAgents);
  return serviceAgents.filter(serviceAgent => serviceAgent.shipment_id === shipmentId);
}
