import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function CreatePpm(moveId, size) {
  const client = await getClient();
  const response = await client.apis.ppm.createPersonallyProcuredMove({
    moveId,
    createPersonallyProcuredMovePayload: { size },
  });
  checkResponse(response, 'failed to create ppm due to server error');
  return response.body;
}
