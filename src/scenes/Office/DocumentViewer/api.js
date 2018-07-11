import { getClient, checkResponse } from 'shared/api';

export async function IndexMoveDocuments(moveId) {
  const client = await getClient();
  const response = await client.apis.moves.indexMoveDocuments({
    moveId,
  });
  checkResponse(response, 'failed to get move documents due to server error');
  return response.body;
}
