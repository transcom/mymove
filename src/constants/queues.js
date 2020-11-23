import { shape, string, bool, arrayOf } from 'prop-types';

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

export const BRANCH_OPTIONS = {
  ARMY: 'Army',
  NAVY: 'Navy',
  MARINES: 'Marine Corps',
  AIR_FORCE: 'Air Force',
  COAST_GUARD: 'Coast Guard',
};

export const SortShape = arrayOf(
  shape({
    id: string,
    desc: bool,
  }),
);
