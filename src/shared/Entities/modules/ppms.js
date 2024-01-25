import { filter } from 'lodash';

import { fetchActivePPM } from '../../utils';

import { swaggerRequest } from 'shared/Swagger/request';
import { getClient } from 'shared/Swagger/api';

const approvePpmLabel = 'PPMs.approvePPM';
export const downloadPPMAttachmentsLabel = 'PPMs.downloadAttachments';
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

export function approveReimbursement(reimbursementId, label = approveReimbursementLabel) {
  const swaggerTag = 'office.approveReimbursement';
  return swaggerRequest(getClient, swaggerTag, { reimbursementId }, { label });
}

export function selectActivePPMForMove(state, moveId) {
  const ppms = Object.values(state.entities.personallyProcuredMoves);
  filter(ppms, (ppm) => ppm.moveId === moveId);
  const activePPM = fetchActivePPM(ppms);
  return activePPM || {};
}
