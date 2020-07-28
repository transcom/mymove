import { getClient, checkResponse } from 'shared/Swagger/api';

export async function GetMove(moveId) {
  const client = await getClient();
  const response = await client.apis.moves.showMove({
    moveId,
  });
  checkResponse(response, 'failed to get move due to server error');
  return response.body;
}

export async function SubmitMoveForApproval(moveId, payload) {
  const client = await getClient();
  const response = await client.apis.moves.submitMoveForApproval({
    moveId,
    submitMoveForApprovalPayload: payload,
  });
  checkResponse(response, 'failed to submit move for approval due to server error');
  return response.body;
}
