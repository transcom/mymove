import { get } from 'lodash';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

const approvePpmLabel = 'PPMs.approvePPM';
export const downloadPPMAttachmentsLabel = 'PPMs.downloadAttachments';

export function approvePPM(personallyProcuredMoveId) {
  const label = approvePpmLabel;
  const swaggerTag = 'office.approvePPM';
  return swaggerRequest(getClient, swaggerTag, { personallyProcuredMoveId }, { label });
}

export function loadPPMs(moveId) {
  const label = 'office.loadPPMs';
  const swaggerTag = 'ppm.indexPersonallyProcuredMoves';
  return swaggerRequest(getClient, swaggerTag, { moveId }, { label });
}

export function updatePPM(
  moveId,
  personallyProcuredMoveId,
  payload /*shape: {size, weightEstimate, estimatedIncentive}*/,
) {
  const label = 'office.updatePPM';
  const swaggerTag = 'ppm.patchPersonallyProcuredMove';
  return swaggerRequest(
    getClient,
    swaggerTag,
    {
      moveId,
      personallyProcuredMoveId,
      patchPersonallyProcuredMovePayload: payload,
    },
    { label },
  );
}

export function downloadPPMAttachments(ppmId, docTypes, label = downloadPPMAttachmentsLabel) {
  const swaggerTag = 'ppm.createPPMAttachments';
  const payload = { personallyProcuredMoveId: ppmId, docTypes };
  return swaggerRequest(getClient, swaggerTag, payload, { label });
}

export function approveReimbursement(reimbursementId) {
  const label = 'office.approveReimbursement';
  const swaggerTag = 'office.approveReimbursement';
  return swaggerRequest(getClient, swaggerTag, { reimbursementId }, { label });
}

export function selectPPMForMove(state, moveId) {
  const ppmId = Object.keys(get(state, 'entities.personallyProcuredMove', {})).find(
    ppmId => get(state, `entities.personallyProcuredMove.${ppmId}.move_id`) === moveId,
  );
  return get(state, `entities.personallyProcuredMove.${ppmId}`, {});
}

export function selectReimbursement(state, advanceId) {
  return get(state, `entities.reimbursements.${advanceId}`);
}
