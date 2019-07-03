import { getClient, checkResponse } from 'shared/Swagger/api';
import { formatPayload, formatDateString } from 'shared/utils';

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
  const response = await client.apis.ppm.patchPersonallyProcuredMove({
    moveId,
    personallyProcuredMoveId,
    patchPersonallyProcuredMovePayload: formatPayload(payload, payloadDef),
  });
  checkResponse(response, 'failed to update ppm due to server error');
  return response.body;
}

export async function GetPpmWeightEstimate(moveDate, originZip, destZip, weightEstimate) {
  const client = await getClient();
  const response = await client.apis.ppm.showPPMEstimate({
    original_move_date: formatDateString(moveDate),
    origin_zip: originZip,
    destination_zip: destZip,
    weight_estimate: weightEstimate,
  });
  checkResponse(response, 'failed to update ppm due to server error');
  return response.body;
}

export async function GetPpmSitEstimate(moveDate, sitDays, originZip, destZip, weightEstimate) {
  const client = await getClient();
  const response = await client.apis.ppm.showPPMSitEstimate({
    original_move_date: formatDateString(moveDate),
    days_in_storage: sitDays,
    origin_zip: originZip,
    destination_zip: destZip,
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
