import { PPM_MAX_ADVANCE_RATIO } from 'constants/shipments';
import { formatCentsTruncateWhole, formatCentsRange, convertCentsToWholeDollarsRoundedDown } from 'utils/formatters';

export const hasShortHaulError = (error) => error?.statusCode === 409;

export const getIncentiveRange = (ppm, estimate) => {
  let range = formatCentsRange(ppm?.incentive_estimate_min, ppm?.incentive_estimate_max);

  if (!range) range = formatCentsRange(estimate?.range_min, estimate?.range_max);

  return range || '';
};

// Calculates the max advance based on the incentive (in cents). Rounds down and returns a cent value as a number.
export const calculateMaxAdvance = (incentive) => {
  return Math.floor(incentive * PPM_MAX_ADVANCE_RATIO);
};

// Calculates max advance and formats max advance and incentive. All values change from cents to dollars and are
// rounded down. Formatted values are strings.
export const calculateMaxAdvanceAndFormatAdvanceAndIncentive = (incentive) => {
  const maxAdvance = calculateMaxAdvance(incentive);

  return {
    maxAdvance: convertCentsToWholeDollarsRoundedDown(maxAdvance),
    formattedMaxAdvance: formatCentsTruncateWhole(maxAdvance),
    formattedIncentive: formatCentsTruncateWhole(incentive),
  };
};

export const getFormattedMaxAdvancePercentage = () => `${PPM_MAX_ADVANCE_RATIO * 100}%`;
