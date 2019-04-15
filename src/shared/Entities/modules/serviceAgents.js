import { swaggerRequest } from 'shared/Swagger/request';
import { getPublicClient } from 'shared/Swagger/api';

const getServiceAgentsForShipmentLabel = 'ServiceAgents.getServiceAgentsForShipment';
const updateServiceAgentForShipmentLabel = 'ServiceAgents.updateServiceAgentForShipment';
const createServiceAgentForShipmentLabel = 'ServiceAgents.createServiceAgentForShipment';

export function getServiceAgentsForShipment(shipmentId, label = getServiceAgentsForShipmentLabel) {
  const swaggerTag = 'service_agents.indexServiceAgents';
  return swaggerRequest(getPublicClient, swaggerTag, { shipmentId }, { label });
}

export function updateServiceAgentForShipment(
  shipmentId,
  serviceAgentId,
  serviceAgent,
  label = updateServiceAgentForShipmentLabel,
) {
  const swaggerTag = 'service_agents.patchServiceAgent';
  return swaggerRequest(getPublicClient, swaggerTag, { shipmentId, serviceAgentId, serviceAgent }, { label });
}

export function createServiceAgentForShipment(shipmentId, serviceAgent, label = createServiceAgentForShipmentLabel) {
  const swaggerTag = 'service_agents.createServiceAgent';
  return swaggerRequest(getPublicClient, swaggerTag, { shipmentId, serviceAgent }, { label });
}

export function updateServiceAgentsForShipment(shipmentId, serviceAgents) {
  return async function(dispatch) {
    Object.values(serviceAgents).map(serviceAgent => {
      if (serviceAgent.id) {
        return dispatch(updateServiceAgentForShipment(shipmentId, serviceAgent.id, serviceAgent));
      } else if (!serviceAgent.company || !serviceAgent.email || !serviceAgent.phone_number) {
        // Don't send the service agent if it's not got enough details
        // Currently, it should only be the destination agent that gets skipped
        return undefined;
      } else {
        return dispatch(createServiceAgentForShipment(shipmentId, serviceAgent));
      }
    });
  };
}

export function selectServiceAgentsForShipment(state, shipmentId) {
  if (!shipmentId) {
    return [];
  }
  const serviceAgents = Object.values(state.entities.serviceAgents);
  return serviceAgents.filter(serviceAgent => serviceAgent.shipment_id === shipmentId);
}
