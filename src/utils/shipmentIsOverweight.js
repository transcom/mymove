// If the shipment's billable weight cap is greater than 110% of the estimated weight,
// then the shipment is overweight
export default function shipmentIsOverweight(estimatedWeight, weightCap) {
  return parseInt(weightCap, 10) > parseInt(estimatedWeight, 10) * 1.1;
}
