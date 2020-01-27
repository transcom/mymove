export function hasShortHaulError(rateEngineError) {
  return rateEngineError && rateEngineError.statusCode === 409 ? true : false;
}
