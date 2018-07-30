import { getClient, checkResponse } from 'shared/api';

// MOVE QUEUE
export async function RetrieveMovesForTSP(queueType) {
  const client = await getClient();
  const response = await client.apis.queues.showQueue({
    queueType,
  });
  checkResponse(response, 'failed to retrieve moves due to server error');
  return response.body;
}
