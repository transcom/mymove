import { getClient, checkResponse } from 'shared/Swagger/api';

export async function CreateMove(ordersId, payload) {
  const client = await getClient();
  const response = await client.apis.moves.createMove({
    ordersId: ordersId,
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

export async function SubmitMoveForApproval(moveId) {
  const client = await getClient();
  const response = await client.apis.moves.submitMoveForApproval({
    moveId,
  });
  checkResponse(
    response,
    'failed to submit move for approval due to server error',
  );
  return response.body;
}
