// PPM Document statuses
export default {
  APPROVED: 'APPROVED',
  REJECTED: 'REJECTED',
  EXCLUDED: 'EXCLUDED',
};

export const PPM_DOCUMENT_STATUS = {
  APPROVED: 'APPROVED',
  REJECTED: 'REJECTED',
  EXCLUDED: 'EXCLUDED',
};

export const ReviewDocumentsStatus = {
  ACCEPT: 'Accept',
  REJECT: 'Reject',
  EXCLUDE: 'Exclude',
};

export const ADVANCE_STATUSES = {
  APPROVED: { apiValue: 'APPROVED', displayValue: 'Approved' },
  REJECTED: { apiValue: 'REJECTED', displayValue: 'Rejected' },
  EDITED: { apiValue: 'EDITED', displayValue: 'Approved' },
  RECEIVED: { apiValue: 'RECEIVED', displayValue: 'Received' },
  NOT_RECEIVED: { apiValue: 'NOT_RECEIVED', displayValue: 'Not received' },
};

export const renderMultiplier = (multiplier) => {
  if (multiplier === '') {
    return null;
  }
  return `(with ${multiplier}x multiplier)`;
};
