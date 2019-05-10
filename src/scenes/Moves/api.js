import { getClient, checkResponse } from 'shared/Swagger/api';

export async function SubmitMoveForApproval(moveId, payload) {
  const client = await getClient();
  const response = await client.apis.moves.submitMoveForApproval({
    moveId,
    submitMoveForApprovalPayload: payload,
  });
  checkResponse(response, 'failed to submit move for approval due to server error');
  return response.body;
}
