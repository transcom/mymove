import { getClient, checkResponse } from 'shared/api';

export async function CreatePpm(moveId, ppmRequest) {
  const client = await getClient();
  const response = await client.apis.ppm.createPersonallyProcuredMove(
    (moveId: moveId),
    (createPersonallyProcuredMovePayload: ppmRequest),
  );
  checkResponse(response, 'failed to create ppm due to server error');
}
