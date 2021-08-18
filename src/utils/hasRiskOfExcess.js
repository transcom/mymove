// If the estimated weight is 90% or more of the weight allowance,
// then there is a risk of excess
export default function hasRiskOfExcess(estimatedWeight, weightAllowance) {
  return 0.9 * weightAllowance <= estimatedWeight;
}
