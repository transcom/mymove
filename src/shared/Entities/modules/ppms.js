import { get } from 'lodash';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

const approvePpmLabel = 'PPMs.approvePPM';
export const downloadPPMAttachmentsLabel = 'PPMs.downloadAttachments';
const loadPPMsLabel = 'office.loadPPMs';
const updatePPMLabel = 'office.updatePPM';
const approveReimbursementLabel = 'office.approveReimbursement';

export function approvePPM(personallyProcuredMoveId, personallyProcuredMoveApproveDate, label = approvePpmLabel) {
  const swaggerTag = 'office.approvePPM';
  return swaggerRequest(
    getClient,
    swaggerTag,
    {
      personallyProcuredMoveId,
      approvePersonallyProcuredMovePayload: {
        approve_date: personallyProcuredMoveApproveDate,
      },
    },
    { label },
  );
}

export function loadPPMs(moveId, label = loadPPMsLabel) {
  const swaggerTag = 'ppm.indexPersonallyProcuredMoves';
  return swaggerRequest(getClient, swaggerTag, { moveId }, { label });
}

export function updatePPM(
  moveId,
  personallyProcuredMoveId,
  payload /*shape: {size, weightEstimate, estimatedIncentive}*/,
  label = updatePPMLabel,
) {
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

export function approveReimbursement(reimbursementId, label = approveReimbursementLabel) {
  const swaggerTag = 'office.approveReimbursement';
  return swaggerRequest(getClient, swaggerTag, { reimbursementId }, { label });
}

export function selectPPMForMove(state, moveId) {
  const ppm = Object.values(state.entities.personallyProcuredMoves).find(ppm => ppm.move_id === moveId);
  return ppm || {};
}

export function selectReimbursement(state, reimbursementId) {
  const advanceFromEntities = get(state, `entities.reimbursements.${reimbursementId}`);
  const advanceFromPpmReducer = get(state, 'ppm.currentPpm.advance');
  return advanceFromEntities || advanceFromPpmReducer || {};
}

export function getMaxAdvance(state, ppmId) {
  const maxIncentive = get(state, `entities.personallyProcuredMoves.${ppmId}.incentive_estimate_max`);
  // we are using 20000000 since it is the largest number MacRae found that could be stored in table
  // and we don't want to block the user from requesting an advance if the rate engine fails
  return maxIncentive ? 0.6 * maxIncentive : 20000000;
}
