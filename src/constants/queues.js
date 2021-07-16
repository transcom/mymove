import { shape, string, bool, arrayOf } from 'prop-types';

import MOVE_STATUSES from 'constants/moves';

export const MOVE_STATUS_OPTIONS = [
  { value: MOVE_STATUSES.SUBMITTED, label: 'New move' },
  { value: MOVE_STATUSES.APPROVALS_REQUESTED, label: 'Approvals requested' },
  { value: MOVE_STATUSES.APPROVED, label: 'Move approved' },
];

// Both moves that progressed straight from customer submission to the TOO
// queue as well as those that completed services counseling should have the
// status label of New move
export const MOVE_STATUS_LABELS = {
  [MOVE_STATUSES.SUBMITTED]: 'New move',
  [MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED]: 'New move',
  [MOVE_STATUSES.APPROVALS_REQUESTED]: 'Approvals requested',
  [MOVE_STATUSES.APPROVED]: 'Move approved',
};

export const SERVICE_COUNSELING_MOVE_STATUS_OPTIONS = [
  { value: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING, label: 'Needs counseling' },
  { value: MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED, label: 'Service counseling completed' },
];

export const SERVICE_COUNSELING_MOVE_STATUS_LABELS = {
  [MOVE_STATUSES.NEEDS_SERVICE_COUNSELING]: 'Needs counseling',
  [MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED]: 'Service counseling completed',
};

export const PAYMENT_REQUEST_STATUS_OPTIONS = [
  { value: 'Payment requested', label: 'Payment requested' },
  { value: 'Reviewed', label: 'Reviewed' },
  { value: 'Rejected', label: 'Rejected' },
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

export const PAGINATION_PAGE_DEFAULT = 1;

export const PAGINATION_PAGE_SIZE_DEFAULT = 20;

export const SortShape = arrayOf(
  shape({
    id: string,
    desc: bool,
  }),
);
