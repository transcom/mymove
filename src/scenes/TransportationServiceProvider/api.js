import { getPublicClient, checkResponse } from 'shared/Swagger/api';

// SHIPMENT QUEUE
export async function RetrieveShipmentsForTSP(queueType) {
  const queueToStatus = {
    new: ['AWARDED'],
    accepted: ['ACCEPTED'],
    approved: ['APPROVED'],
    in_transit: ['IN_TRANSIT'],
    delivered: ['DELIVERED'],
    completed: ['COMPLETED'],
    all: [],
  };
  /* eslint-disable security/detect-object-injection */
  const status = (queueType && queueToStatus[queueType] && queueToStatus[queueType].join(',')) || '';
  /* eslint-enable security/detect-object-injection */
  const client = await getPublicClient();
  const response = await client.apis.shipments.indexShipments({
    status,
  });
  checkResponse(response, 'failed to retrieve moves due to server error');
  return response.body;
}

// All documents for shipment
export async function GetAllShipmentDocuments(shipmentId) {
  const client = await getPublicClient();
  const response = await client.apis.move_docs.indexMoveDocuments({
    shipmentId,
  });
  checkResponse(response, 'failed to load shipment documents due to server error');
  return response.body;
}
