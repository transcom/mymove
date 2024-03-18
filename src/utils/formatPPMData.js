/* eslint-disable camelcase */
import { formatMoveHistoryFullAddressFromJSON } from './formatters';

import { PPM_UPLOAD_TYPES, PPM_UPLOAD_TYPES_LABELS } from 'constants/ppmUploadTypes';

const formatDisplayName = (uploadType) => {
  if (PPM_UPLOAD_TYPES.some((type) => type === uploadType)) {
    return PPM_UPLOAD_TYPES_LABELS[uploadType];
  }

  return '';
};

export const formatCloseoutOfficeFor = ({ changedValues = {}, context: [{ closeout_office_name }] = [{}] }) =>
  'closeout_office_id' in changedValues &&
  closeout_office_name && { closeout_office_name: changedValues.closeout_office_id && closeout_office_name };

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

  return ppmValues;
}

export default { formatDataForPPM };
