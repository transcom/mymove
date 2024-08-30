import { shape, string, bool, arrayOf } from 'prop-types';

import { roleTypes } from './userRoles';

import MOVE_STATUSES from 'constants/moves';

export const MOVE_STATUS_OPTIONS = [
  { value: MOVE_STATUSES.SUBMITTED, label: 'New move' },
  { value: MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED, label: 'Service Counseling Completed' },
  { value: MOVE_STATUSES.APPROVALS_REQUESTED, label: 'Approvals requested' },
];

// Both moves that progressed straight from customer submission to the TOO
// queue as well as those that completed services counseling should have the
// status label of New move
export const MOVE_STATUS_LABELS = {
  [MOVE_STATUSES.DRAFT]: 'Draft',
  [MOVE_STATUSES.SUBMITTED]: 'New move',
  [MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED]: 'Service Counseling Completed',
  [MOVE_STATUSES.NEEDS_SERVICE_COUNSELING]: 'Needs Service Counseling',
  [MOVE_STATUSES.APPROVALS_REQUESTED]: 'Approvals requested',
  [MOVE_STATUSES.APPROVED]: 'Move approved',
};

export const SEARCH_QUEUE_STATUS_FILTER_OPTIONS = [
  { value: MOVE_STATUSES.DRAFT, label: 'Draft' },
  { value: MOVE_STATUSES.SUBMITTED, label: 'New Move' },
  { value: MOVE_STATUSES.NEEDS_SERVICE_COUNSELING, label: 'Needs counseling' },
  { value: MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED, label: 'Service counseling completed' },
  { value: MOVE_STATUSES.APPROVED, label: 'Move Approved' },
];

export const SERVICE_COUNSELING_MOVE_STATUS_LABELS = {
  [MOVE_STATUSES.NEEDS_SERVICE_COUNSELING]: 'Needs counseling',
  [MOVE_STATUSES.SERVICE_COUNSELING_COMPLETED]: 'Service counseling completed',
};

export const PAYMENT_REQUEST_STATUS_OPTIONS = [{ value: 'PENDING', label: 'Payment requested' }];

export const ROLE_TYPE_OPTIONS = {
  [roleTypes.SERVICES_COUNSELOR]: SEARCH_QUEUE_STATUS_FILTER_OPTIONS,
  [roleTypes.CUSTOMER_SERVICE_REPRESENTATIVE]: MOVE_STATUS_OPTIONS,
  [roleTypes.QAE]: MOVE_STATUS_OPTIONS,
  [roleTypes.TOO]: MOVE_STATUS_OPTIONS,
  [roleTypes.TIO]: PAYMENT_REQUEST_STATUS_OPTIONS,
};

export const BRANCH_OPTIONS = [
  { value: '', label: 'All' },
  { value: 'ARMY', label: 'Army' },
  { value: 'NAVY', label: 'Navy' },
  { value: 'AIR_FORCE', label: 'Air Force' },
  { value: 'COAST_GUARD', label: 'Coast Guard' },
  { value: 'SPACE_FORCE', label: 'Space Force' },
  { value: 'MARINES', label: 'Marine Corps' },
];

export const SERVICE_COUNSELING_PPM_STATUS_OPTIONS = [
  { value: '', label: 'All' },
  { value: 'WAITING_ON_CUSTOMER', label: 'Waiting on customer' },
  { value: 'NEEDS_CLOSEOUT', label: 'Needs closeout' },
];

export const SERVICE_COUNSELING_PPM_STATUS_LABELS = {
  CANCELLED: 'Cancelled',
  DRAFT: 'Draft',
  SUBMITTED: 'Submitted',
  WAITING_ON_CUSTOMER: 'Waiting on customer',
  NEEDS_ADVANCE_APPROVAL: 'Needs advance approval',
  NEEDS_CLOSEOUT: 'Needs closeout',
  CLOSEOUT_COMPLETE: 'Closeout complete',
  COMPLETED: 'Completed',
};

export const SERVICE_COUNSELING_PPM_TYPE_OPTIONS = [
  { value: '', label: 'All' },
  { value: 'FULL', label: 'Full' },
  { value: 'PARTIAL', label: 'Partial' },
];

export const SERVICE_COUNSELING_PPM_TYPE_LABELS = {
  FULL: 'Full',
  PARTIAL: 'Partial',
};

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
