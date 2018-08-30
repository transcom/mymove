import { getPublicClient, checkResponse } from 'shared/api';
import { formatPayload } from 'shared/utils';

// SHIPMENT QUEUE
export async function RetrieveShipmentsForTSP(queueType) {
  const queueToStatus = {
    new: ['AWARDED'],
    all: [],
  };
  /* eslint-disable security/detect-object-injection */
  const status =
    (queueType &&
      queueToStatus[queueType] &&
      queueToStatus[queueType].join(',')) ||
    '';
  /* eslint-enable security/detect-object-injection */
  const client = await getPublicClient();
  const response = await client.apis.shipments.indexShipments({
    status,
    limit: 25,
    offset: 1,
  });
  checkResponse(response, 'failed to retrieve moves due to server error');
  return response.body;
}

// SHIPMENT
export async function LoadShipment(shipmentId) {
  const client = await getPublicClient();
  const response = await client.apis.shipments.getShipment({
    shipmentId,
  });
  checkResponse(response, 'failed to load shipment due to server error');
  return response.body;
}

// SHIPMENT ACCEPT
export async function AcceptShipment(shipmentId) {
  const client = await getPublicClient();
  const response = await client.apis.shipments.createShipmentAccept({
    shipmentId: shipmentId,
  });
  checkResponse(response, 'failed to accept shipment due to server error');
  return response.body;
}

export async function PatchShipment(shipmentId, shipment) {
  const client = await getPublicClient();
  const payloadDef = client.spec.definitions.Shipment;
  const response = await client.apis.shipments.patchShipment({
    shipmentId,
    update: formatPayload(shipment, payloadDef),
  });
  checkResponse(response, 'failed to load shipment due to server error');
  return response.body;
}

// ServiceAgents
export async function CreateServiceAgent(shipmentId, payload) {
  const client = await getPublicClient();
  const response = await client.apis.service_agents.createServiceAgent({
    shipment_uuid: shipmentId,
    payload,
  });
  checkResponse(response, 'failed to load service agent due to server error');
  return response.body;
}

export async function UpdateServiceAgent(payload) {
  const client = await getPublicClient();
  const response = await client.apis.service_agents.updateServiceAgents({
    shipment_uuid: payload.id,
    payload,
  });
  checkResponse(response, 'failed to update service agent due to server error');
  return response.body;
}

export async function IndexServiceAgents(shipmentId) {
  const client = await getPublicClient();
  const response = await client.apis.service_agents.indexServiceAgents({
    shipment_uuid: shipmentId,
  });
  checkResponse(response, 'failed to load service agents due to server error');
  return response.body;
}
