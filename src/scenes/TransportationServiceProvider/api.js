import { getPublicClient, checkResponse } from 'shared/Swagger/api';
import { formatPayload } from 'shared/utils';

// SHIPMENT QUEUE
export async function RetrieveShipmentsForTSP(queueType) {
  const queueToStatus = {
    new: ['AWARDED'],
    in_transit: ['IN_TRANSIT'],
    approved: ['APPROVED'],
    delivered: ['DELIVERED'],
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
  const response = await client.apis.shipments.acceptShipment({
    shipmentId: shipmentId,
  });
  checkResponse(response, 'failed to accept shipment due to server error');
  return response.body;
}

export async function RejectShipment(shipmentId, reason) {
  const client = await getPublicClient();
  const response = await client.apis.shipments.rejectShipment({
    shipmentId: shipmentId,
    payload: { reason },
  });
  checkResponse(response, 'failed to accept shipment due to server error');
  return response.body;
}

export async function TransportShipment(shipmentId, payload) {
  const client = await getPublicClient();
  const payloadDef = client.spec.definitions.ActualPickupDate;
  const response = await client.apis.shipments.transportShipment({
    shipmentId,
    payload: formatPayload(payload, payloadDef),
  });
  checkResponse(response, 'failed to pick up shipment due to server error');
  return response.body;
}

export async function DeliverShipment(shipmentId, payload) {
  const client = await getPublicClient();
  const payloadDef = client.spec.definitions.ActualDeliveryDate;
  const response = await client.apis.shipments.deliverShipment({
    shipmentId,
    payload: formatPayload(payload, payloadDef),
  });
  checkResponse(response, 'failed to pick up shipment due to server error');
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
    shipmentId,
    serviceAgent: payload,
  });
  checkResponse(response, 'failed to create service agent due to server error');
  return response.body;
}

export async function UpdateServiceAgent(payload) {
  const client = await getPublicClient();
  const response = await client.apis.service_agents.patchServiceAgent({
    shipmentId: payload.shipment_id,
    serviceAgentId: payload.id,
    patchServiceAgentPayload: payload,
  });
  checkResponse(response, 'failed to update service agent due to server error');
  return response.body;
}

export async function IndexServiceAgents(shipmentId) {
  const client = await getPublicClient();
  const response = await client.apis.service_agents.indexServiceAgents({
    shipmentId,
  });
  checkResponse(response, 'failed to load service agents due to server error');
  return response.body;
}

// Generate Gov Bill of Lading
export async function GenerateGBL(shipmentId) {
  const client = await getPublicClient();
  const response = await client.apis.shipments.createGovBillOfLading({
    shipmentId: shipmentId,
  });
  checkResponse(
    response,
    'failed to create government bill of lading due to server error',
  );
  return response.body;
}

// All documents for shipment
export async function GetAllShipmentDocuments(shipmentId) {
  const client = await getPublicClient();
  const response = await client.apis.move_docs.indexMoveDocuments({
    shipmentId,
  });
  checkResponse(
    response,
    'failed to load shipment documents due to server error',
  );
  return response.body;
}
