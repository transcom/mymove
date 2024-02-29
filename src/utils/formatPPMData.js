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

  return ppmValues;
}

export default { formatDataForPPM };
