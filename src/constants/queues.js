// eslint-disable-next-line import/prefer-default-export
export const MOVE_STATUS_OPTIONS = {
  SUBMITTED: 'New move',
  'APPROVALS REQUESTED': 'Approvals requested',
  APPROVED: 'Move approved',
};

export const PAYMENT_REQUEST_STATUS_OPTIONS = {
  'Payment requested': 'Payment requested',
  Reviewed: 'Reviewed',
  Paid: 'Paid',
};

export const BRANCH_OPTIONS = [
  { value: '', label: 'All' },
  { value: 'ARMY', label: 'Army' },
  { value: 'NAVY', label: 'Navy' },
  { value: 'MARINES', label: 'Marine Corps' },
  { value: 'AIR_FORCE', label: 'Air Force' },
  { value: 'COAST_GUARD', label: 'Coast Guard' },
];

export const BRANCH_OPTIONS_NO_MARINES = [
  { value: '', label: 'All' },
  { value: 'ARMY', label: 'Army' },
  { value: 'NAVY', label: 'Navy' },
  { value: 'AIR_FORCE', label: 'Air Force' },
  { value: 'COAST_GUARD', label: 'Coast Guard' },
];

export const GBLOC = {
  USMC: 'USMC',
};
