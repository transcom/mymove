// If the shipment's billable weight cap is greater than 110% of the estimated weight,
// then the shipment is overweight
// eslint-disable-next-line import/prefer-default-export
export function shipmentIsOverweight(estimatedWeight, weightCap) {
  return weightCap > estimatedWeight * 1.1;
}
