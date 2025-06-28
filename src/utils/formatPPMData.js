/* eslint-disable camelcase */
import { formatMoveHistoryFullAddressFromJSON } from './formatters';

import { PPM_UPLOAD_TYPES, PPM_UPLOAD_TYPES_LABELS } from 'constants/ppmUploadTypes';

const formatDisplayName = (uploadType) => {
  if (PPM_UPLOAD_TYPES.some((type) => type === uploadType)) {
    return PPM_UPLOAD_TYPES_LABELS[uploadType];
  }

  return '';
};

export const formatCloseoutOfficeFor = ({ changedValues = {}, context = [{}] }) =>
  'closeout_office_id' in changedValues &&
  context[0]?.closeout_office_name && {
    closeout_office_name: changedValues.closeout_office_id && context[0].closeout_office_name,
  };

export const formatUploadTypeFor = ({ context: [{ upload_type }] = [{}] }) =>
  upload_type && { upload_type: formatDisplayName(upload_type) };

export const formatFileNameFor = ({ context: [{ filename }] = [{}] }) => filename && { filename };

export const formatW2AddressFor = ({ changedValues: { w2_address_id } = {}, context: [{ w2_address }] = [{}] }) =>
  w2_address_id && w2_address && { w2_address: formatMoveHistoryFullAddressFromJSON(w2_address) };

export function formatDataForPPM(historyRecord) {
  const ppmValues = {
    ...formatCloseoutOfficeFor(historyRecord),
    ...formatUploadTypeFor(historyRecord),
    ...formatFileNameFor(historyRecord),
    ...formatW2AddressFor(historyRecord),
  };

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
