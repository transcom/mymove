import { get, isNull } from 'lodash';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
import { loadEntitlementsFromState } from 'shared/entitlements';
import { getCurrentShipmentID, getCurrentMove } from 'shared/UI/ducks';
import { selectShipment } from 'shared/Entities/modules/shipments';

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

export function isHHGPPMComboMove(state) {
  const move = getCurrentMove(state);
  return get(move, 'selected_move_type') === 'HHG_PPM';
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

const estimatedRemainingWeight = (sum, weight) => {
  if (sum >= weight) {
    return sum - weight;
  } else {
    return sum;
  }
};

export function getEstimatedRemainingWeight(state) {
  const entitlements = loadEntitlementsFromState(state);

  if (!isHHGPPMComboMove(state) || isNull(entitlements)) {
    return null;
  }

  const { sum } = entitlements;

  const { pm_survey_weight_estimate, weight_estimate } = selectShipment(state, getCurrentShipmentID(state));

  if (pm_survey_weight_estimate) {
    return estimatedRemainingWeight(sum, pm_survey_weight_estimate);
  }

  if (sum && weight_estimate >= 0) {
    return estimatedRemainingWeight(sum, weight_estimate);
  }
}

export function getActualRemainingWeight(state) {
  const entitlements = loadEntitlementsFromState(state);

  if (!isHHGPPMComboMove(state) || isNull(entitlements)) {
    return null;
  }

  const { sum } = entitlements;
  const { tare_weight, gross_weight } = selectShipment(state, getCurrentShipmentID(state));

  if (sum && gross_weight && tare_weight) {
    return estimatedRemainingWeight(sum, gross_weight - tare_weight);
  }
}
