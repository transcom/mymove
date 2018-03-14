import { getClient, checkResponse } from 'shared/api';

export async function CreateMove(selectedMoveType) {
  const client = await getClient();
  const response = await client.apis.moves.createMove({
    createMovePayload: selectedMoveType,
  });
  checkResponse(response, 'failed to create move due to server error');
  return response.body;
}
