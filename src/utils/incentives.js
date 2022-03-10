import { formatCentsTruncateWhole } from './formatters';

import { formatCentsRange } from 'shared/formatters';

export const hasShortHaulError = (error) => error?.statusCode === 409;

export const getIncentiveRange = (ppm, estimate) => {
  let range = formatCentsRange(ppm?.incentive_estimate_min, ppm?.incentive_estimate_max);

  if (!range) range = formatCentsRange(estimate?.range_min, estimate?.range_max);

  return range || '';
};

// returns 60% of the incentive in dollars, rounded down to nearest whole number
export const maxAdvance = (incentive) => {
  // incentive is in cents, convert to dollars rounded down to nearest whole number
  const incentiveInDollars = formatCentsTruncateWhole(incentive);
  // max advance is equal to 60% of the incentive, rounded down to the nearest whole number
  return Math.floor(incentiveInDollars * 0.6);
};
