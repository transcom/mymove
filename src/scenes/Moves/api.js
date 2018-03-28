import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function CreateMove(payload) {
  const client = await getClient();
  const response = await client.apis.moves.createMove({
    createMovePayload: payload,
  });
  checkResponse(response, 'failed to create move due to server error');
  return response.body;
}

export async function GetMove(moveId) {
  const client = await getClient();
  const response = await client.apis.moves.showMove({
    moveId,
  });
  checkResponse(response, 'failed to get move due to server error');
  return response.body;
}

export async function UpdateMove(
  moveId,
  payload /*shape: { selected_move_type }*/,
) {
  const client = await getClient();
  const response = await client.apis.moves.patchMove({
    moveId,
    patchMovePayload: payload,
  });
  checkResponse(response, 'failed to update move due to server error');
  return response.body;
}
