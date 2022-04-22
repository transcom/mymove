import { formatCentsTruncateWhole, formatCentsRange } from './formatters';

export const hasShortHaulError = (error) => error?.statusCode === 409;

export const getIncentiveRange = (ppm, estimate) => {
  let range = formatCentsRange(ppm?.incentive_estimate_min, ppm?.incentive_estimate_max);

  if (!range) range = formatCentsRange(estimate?.range_min, estimate?.range_max);

  return range || '';
};

// MaxAdvance returns 60% of the incentive in dollars, rounded down to nearest whole number
// As a formated string
export const maxAdvance = (incentive) => {
  return formatCentsTruncateWhole(Math.floor(incentive * 0.6));
};
