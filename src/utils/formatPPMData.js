import { formatMoveHistoryFullAddressFromJSON } from './formatters';

import { PPM_UPLOAD_TYPES, PPM_UPLOAD_TYPES_LABELS } from 'constants/ppmUploadTypes';

const formatDisplayName = (uploadType) => {
  if (PPM_UPLOAD_TYPES.some((type) => type === uploadType)) {
    return PPM_UPLOAD_TYPES_LABELS[uploadType];
  }

  return '';
};

export function formatDataForPPM(historyRecord) {
  const ppmValues = {};

  if (historyRecord.context[0].upload_type) {
    ppmValues.upload_type = formatDisplayName(historyRecord.context[0].upload_type);
  }

  if (historyRecord.context[0]?.filename) {
    ppmValues.filename = historyRecord.context[0]?.filename;
  }

  if (historyRecord.changedValues.w2_address_id) {
    ppmValues.w2_address = formatMoveHistoryFullAddressFromJSON(historyRecord.context[0].w2_address);
  }

  // it was requested that we add a status of 'ADDED' once a customer finishes creating a ppm
  // will require refactor down the line but this has been approved for now
  if (
    historyRecord.changedValues?.has_requested_advance !== undefined &&
    historyRecord.oldValues?.has_requested_advance === null &&
    historyRecord.oldValues?.advance_amount_requested === null
  ) {
    ppmValues.ppm_status = 'ADDED';
  }

  return ppmValues;
}

export default { formatDataForPPM };
