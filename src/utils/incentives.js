import { formatCentsRange } from 'shared/formatters';

export const hasShortHaulError = (error) => error?.statusCode === 409;

export const getIncentiveRange = (ppm, estimate) => {
  let range = formatCentsRange(ppm?.incentive_estimate_min, ppm?.incentive_estimate_max);

  if (!range) range = formatCentsRange(estimate?.range_min, estimate?.range_max);

  return range || '';
};

export const maxAdvance = (incentive) => {
  return Math.floor(incentive * 0.6);
};
