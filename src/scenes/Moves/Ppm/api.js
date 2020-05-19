import { getClient, checkResponse } from 'shared/Swagger/api';
import { formatPayload } from 'shared/utils';
import { formatDateForSwagger } from 'shared/dates';

export async function GetPpm(moveId) {
  const client = await getClient();
  const response = await client.apis.ppm.indexPersonallyProcuredMoves({
    moveId,
  });
  checkResponse(response, 'failed to get ppm due to server error');
  return response.body;
}

export async function CreatePpm(moveId, payload /*shape: {size, weightEstimate, estimatedIncentive}*/) {
  const client = await getClient();
  const payloadDef = client.spec.definitions.CreatePersonallyProcuredMovePayload;
  const response = await client.apis.ppm.createPersonallyProcuredMove({
    moveId,
    createPersonallyProcuredMovePayload: formatPayload(payload, payloadDef),
  });
  checkResponse(response, 'failed to create ppm due to server error');
  return response.body;
}

export async function UpdatePpm(
  moveId,
  personallyProcuredMoveId,
  payload /*shape: {size, weightEstimate, estimatedIncentive}*/,
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

export async function GetPpmSitEstimate(moveDate, sitDays, originZip, ordersID, weightEstimate) {
  const client = await getClient();
  const response = await client.apis.ppm.showPPMSitEstimate({
    original_move_date: formatDateForSwagger(moveDate),
    days_in_storage: sitDays,
    origin_zip: originZip,
    orders_id: ordersID,
    weight_estimate: weightEstimate,
  });
  checkResponse(response, 'failed to update ppm due to server error');
  return response.body;
}

export async function RequestPayment(personallyProcuredMoveId) {
  const client = await getClient();
  const response = await client.apis.ppm.requestPPMPayment({
    personallyProcuredMoveId,
  });
  checkResponse(response, 'failed to update ppm status due to server error');
  return response.body;
}
