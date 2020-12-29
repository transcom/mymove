import { get } from 'lodash';
import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';
import { formatDateForSwagger } from 'shared/dates';

const approvePpmLabel = 'PPMs.approvePPM';
export const downloadPPMAttachmentsLabel = 'PPMs.downloadAttachments';
const loadPPMsLabel = 'office.loadPPMs';
const createPPMLabel = 'office.createPPM';
const updatePPMLabel = 'office.updatePPM';
const updatePPMEstimateLabel = 'ppm.updatePPMEstimate';
const approveReimbursementLabel = 'office.approveReimbursement';
const getPPMEstimateLabel = 'ppm.showPPMEstimate';
const getPPMSitEstimateLabel = 'ppm.updatePPMEstimate';

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

export function createPPM(
  moveId,
  payload /*shape: {size, weightEstimate, estimatedIncentive}*/,
  label = createPPMLabel,
) {
  const swaggerTag = 'ppm.createPersonallyProcuredMove';
  payload.original_move_date = formatDateForSwagger(payload.original_move_date);
  payload.actual_move_date = formatDateForSwagger(payload.actual_move_date);
  return swaggerRequest(
    getClient,
    swaggerTag,
    {
      moveId,
      createPersonallyProcuredMovePayload: payload,
    },
    { label },
  );
}

export function updatePPM(
  moveId,
  personallyProcuredMoveId,
  payload /*shape: {size, weightEstimate, estimatedIncentive}*/,
  label = updatePPMLabel,
) {
  const swaggerTag = 'ppm.patchPersonallyProcuredMove';
  payload.original_move_date = formatDateForSwagger(payload.original_move_date);
  payload.actual_move_date = formatDateForSwagger(payload.actual_move_date);
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

export function getPpmWeightEstimate(
  moveDate,
  originZip,
  originDutyStationZip,
  ordersId,
  weightEstimate,
  label = getPPMEstimateLabel,
) {
  const swaggerTag = 'ppm.showPPMEstimate';
  const schemaKey = 'ppmEstimateRange';
  return swaggerRequest(
    getClient,
    swaggerTag,
    {
      original_move_date: moveDate,
      origin_zip: originZip,
      origin_duty_station_zip: originDutyStationZip,
      orders_id: ordersId,
      weight_estimate: weightEstimate,
    },
    { label, schemaKey },
  );
}

export function updatePPMEstimate(moveId, personallyProcuredMoveId, label = updatePPMEstimateLabel) {
  const swaggerTag = 'ppm.updatePersonallyProcuredMoveEstimate';
  return swaggerRequest(
    getClient,
    swaggerTag,
    {
      moveId,
      personallyProcuredMoveId,
    },
    {
      label,
    },
  );
}

export function getPPMSitEstimate(
  ppmId,
  moveDate,
  sitDays,
  originZip,
  ordersID,
  weightEstimate,
  label = getPPMSitEstimateLabel,
) {
  const swaggerTag = 'ppm.showPPMSitEstimate';
  const schemaKey = 'ppmSitEstimate';
  const payload = {
    personally_procured_move_id: ppmId,
    original_move_date: formatDateForSwagger(moveDate),
    days_in_storage: sitDays,
    origin_zip: originZip,
    orders_id: ordersID,
    weight_estimate: weightEstimate,
  };
  return swaggerRequest(getClient, swaggerTag, payload, { label, schemaKey });
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

export function selectPPMEstimateRange(state) {
  if (state.entities.ppmEstimateRanges) {
    return state.entities.ppmEstimateRanges.undefined;
  }
  return {};
}

export function selectPPMSitEstimate(state) {
  if (state.entities.ppmSitEstimate) {
    return state.entities.ppmSitEstimate.undefined.estimate;
  }
  return '';
}

export function selectReimbursement(state, reimbursementId) {
  const advanceFromEntities = get(state, `entities.reimbursements.${reimbursementId}`);
  // todo
  const advanceFromPpmReducer = get(state, 'ppm.currentPpm.advance');
  return advanceFromEntities || advanceFromPpmReducer || {};
}
