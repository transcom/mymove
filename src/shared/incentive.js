import { formatCentsRange } from 'shared/formatters';
import { isEmpty } from 'lodash';

export function hasShortHaulError(rateEngineError) {
  return rateEngineError && rateEngineError.statusCode === 409 ? true : false;
}

export function formatIncentiveRange(ppm, ppmEstimateRange) {
  let incentiveRange = formatCentsRange(ppm.incentive_estimate_min, ppm.incentive_estimate_max);

  // workaround for incentive not yet saved on entities ppm
  if (incentiveRange === '' && !isEmpty(ppmEstimateRange)) {
    incentiveRange = formatCentsRange(ppmEstimateRange.range_min, ppmEstimateRange.range_max);
  }
  // work around for for ppm redux storage in multiple places...
  if (incentiveRange === '' && ppm.currentPpm) {
    incentiveRange = formatCentsRange(ppm.currentPpm.incentive_estimate_min, ppm.currentPpm.incentive_estimate_max);
  }
  return incentiveRange;
}
