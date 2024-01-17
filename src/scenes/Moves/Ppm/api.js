import { getClient, checkResponse } from 'shared/Swagger/api';
import { formatPayload } from 'shared/utils';
import { formatDateForSwagger } from 'shared/dates';

export async function UpdatePpm(
  moveId,
  personallyProcuredMoveId,
  payload /* shape: {size, weightEstimate, estimatedIncentive} */,
) {
  const client = await getClient();
  const payloadDef = client.spec.definitions.PatchPersonallyProcuredMovePayload;
  payload.original_move_date = formatDateForSwagger(payload.original_move_date);
  const response = await client.apis.ppm.patchPersonallyProcuredMove({
    moveId,
    personallyProcuredMoveId,
    patchPersonallyProcuredMovePayload: formatPayload(payload, payloadDef),
  });
  checkResponse(response, 'failed to update ppm due to server error');
  return response.body;
}
