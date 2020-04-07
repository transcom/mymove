import { formatCentsRange } from 'shared/formatters';

export function hasShortHaulError(rateEngineError) {
  return rateEngineError && rateEngineError.statusCode === 409 ? true : false;
}

export function formatIncentiveRange(ppm) {
  let incentiveRange = formatCentsRange(ppm.incentive_estimate_min, ppm.incentive_estimate_max);
  // work around for for ppm redux storage in multiple places...
  if (incentiveRange === '' && ppm.currentPpm) {
    incentiveRange = formatCentsRange(ppm.currentPpm.incentive_estimate_min, ppm.currentPpm.incentive_estimate_max);
  }
  return incentiveRange;
}
