import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function CreateMove(selectedMoveType) {
  const client = await getClient();
  console.log('******', selectedMoveType);
  const response = await client.apis.moves.createMove({
    createMovePayload: selectedMoveType,
  });
  checkResponse(response, 'failed to create move due to server error');
  return response.body;
}

export async function UpdateMove(
  moveId,
  payload /*shape: { selected_move_type }*/,
) {
  const client = await getClient();
  console.log('******', moveId, payload);
  //todo: this is failing due to "Unhandled Rejection (Error): Required parameter personallyProcuredMovePayload is not provided"
  const response = await client.apis.moves.patchMove({
    moveId,
    patchMovePayload: payload,
  });
  checkResponse(response, 'failed to update move due to server error');
  return response.body;
}
