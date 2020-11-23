// eslint-disable-next-line import/prefer-default-export
export const MOVE_STATUS_OPTIONS = [
  { value: 'SUBMITTED', label: 'New move' },
  { value: 'APPROVALS REQUESTED', label: 'Approvals requested' },
  { value: 'APPROVED', label: 'Move approved' },
];

export const PAYMENT_REQUEST_STATUS_OPTIONS = [
  { value: 'Payment requested', label: 'Payment requested' },
  { value: 'Reviewed', label: 'Reviewed' },
  { value: 'Paid', label: 'Paid' },
];

export const BRANCH_OPTIONS = [
  { value: '', label: 'All' },
  { value: 'ARMY', label: 'Army' },
  { value: 'NAVY', label: 'Navy' },
  { value: 'AIR_FORCE', label: 'Air Force' },
  { value: 'COAST_GUARD', label: 'Coast Guard' },
];

export const GBLOC = {
  USMC: 'USMC',
};
