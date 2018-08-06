import { getPublicClient, checkResponse } from 'shared/api';

// MOVE QUEUE
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
