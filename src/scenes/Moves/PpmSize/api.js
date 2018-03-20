import { getClient, checkResponse } from 'shared/api';

export async function GetSpec() {
  const client = await getClient();
  return client.spec;
}

export async function CreatePpm(ppmRequest) {
  const client = await getClient();
  const response = await client.apis.ppm.createPersonallyProcuredMove({
    createPersonallyProcuredMovePayload: ppmRequest,
  });
  checkResponse(response, 'failed to create ppm due to server error');
  return response.body;
}
