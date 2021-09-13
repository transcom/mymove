// If the shipment's billable weight cap is greater than 110% of the estimated weight,
// then the shipment is overweight
export function shipmentIsOverweight(estimatedWeight, weightCap) {
  return weightCap > estimatedWeight * 1.1;
}

export function calcWeightRequested(mtoShipments) {
  return mtoShipments.reduce((accum, shipment) => {
    if (!shipment.reweigh?.weight) {
      return accum + shipment.primeActualWeight;
    }
    return accum + Math.min(shipment.primeActualWeight, shipment.reweigh.weight);
  }, 0);
}

export function calcTotalBillableWeight(mtoShipments) {
  return mtoShipments.reduce((accum, shipment) => {
    if (shipment.billableWeightCap) {
      return accum + shipment.billableWeightCap;
    }

    if (!shipment.reweigh?.weight) {
      return accum + shipment.primeActualWeight;
    }

    return accum + Math.min(shipment.primeActualWeight, shipment.reweigh.weight);
  }, 0);
}

export function calcTotalEstimatedWeight(mtoShipments) {
  return mtoShipments.reduce((accum, shipment) => accum + shipment.primeEstimatedWeight, 0);
}
