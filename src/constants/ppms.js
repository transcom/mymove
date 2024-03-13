// PPM Document statuses
export default {
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
  EDITED: { apiValue: 'EDITED', displayValue: 'Approved with adjustment' },
};
