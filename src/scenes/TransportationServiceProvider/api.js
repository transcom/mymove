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
    shipment_uuid: shipmentId,
  });
  checkResponse(response, 'failed to load shipment due to server error');
  return response.body;
}

export async function PatchShipment(shipmentId, shipment) {
  const client = await getPublicClient();
  const payloadDef = client.spec.definitions.Shipment;
  const response = await client.apis.shipments.patchShipment({
    shipment_uuid: shipmentId,
    update: formatPayload(shipment, payloadDef),
  });
  checkResponse(response, 'failed to load shipment due to server error');
  return response.body;
}

// ServiceAgents
export async function CreateServiceAgent(shipmentId, payload) {
  const client = await getPublicClient();
  const response = await client.apis.shipments.createServiceAgent({
    shipment_uuid: shipmentId,
    payload,
  });
  checkResponse(response, 'failed to load shipment due to server error');
  return response.body;
}

export async function IndexServiceAgents(shipmentId) {
  const client = await getPublicClient();
  const response = await client.apis.shipments.indexServiceAgent({
    shipment_uuid: shipmentId,
  });
  checkResponse(response, 'failed to load shipment due to server error');
  return response.body;
}
