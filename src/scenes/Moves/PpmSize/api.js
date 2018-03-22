import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}
export async function GetPpm(moveId) {
  const client = await getClient();
  const response = await client.apis.ppm.indexPersonallyProcuredMoves({
    moveId,
  });
  checkResponse(response, 'failed to create ppm due to server error');
  return response.body;
}

export async function CreatePpm(
  moveId,
  payload /*shape: {size, weightEstimate}*/,
) {
  const client = await getClient();
  const response = await client.apis.ppm.createPersonallyProcuredMove({
    moveId,
    createPersonallyProcuredMovePayload: payload,
  });
  checkResponse(response, 'failed to create ppm due to server error');
  return response.body;
}

export async function UpdatePpm(
  moveId,
  personallyProcuredMoveId,
  payload /*shape: {size, weightEstimate}*/,
) {
  const client = await getClient();
  console.log('******', personallyProcuredMoveId);
  //todo: this is failing due to "Unhandled Rejection (Error): Required parameter personallyProcuredMovePayload is not provided"
  const response = await client.apis.ppm.updatePersonallyProcuredMove({
    moveId,
    personallyProcuredMoveId,
    createPersonallyProcuredMovePayload: payload,
  });
  checkResponse(response, 'failed to update ppm due to server error');
  return response.body;
}
