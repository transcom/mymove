// If the shipment's billable weight cap is greater than 110% of the estimated weight,
// then the shipment is overweight

import returnLowestValue from './returnLowestValue';

// eslint-disable-next-line import/prefer-default-export
export function shipmentIsOverweight(estimatedWeight, weightCap) {
  return weightCap > estimatedWeight * 1.1;
}

export const getShipmentEstimatedWeight = (shipment) => {
  if (shipment.ppmShipment) {
    return shipment.ppmShipment.estimatedWeight ?? 0;
  }
  return shipment.primeEstimatedWeight ?? 0;
};

export const calculateNetWeightForProGearWeightTicket = (weightTicket) => {
  if (weightTicket.weight == null || Number.isNaN(Number(weightTicket.weight))) {
    return 0;
  }

  return weightTicket.weight;
};

export const calculateTotalNetWeightForProGearWeightTickets = (proGearWeightTickets = []) => {
  return proGearWeightTickets.reduce((prev, curr) => {
    return prev + calculateNetWeightForProGearWeightTicket(curr);
  }, 0);
};

export const calculateWeightTicketWeightDifference = (weightTicket) => {
  if (
    weightTicket.emptyWeight == null ||
    weightTicket.fullWeight == null ||
    Number.isNaN(Number(weightTicket.emptyWeight)) ||
    Number.isNaN(Number(weightTicket.fullWeight))
  ) {
    return 0;
  }

  return weightTicket.fullWeight - weightTicket.emptyWeight;
};

export const getWeightTicketNetWeight = (weightTicket) => {
  if (weightTicket.status !== 'REJECTED')
    return weightTicket.adjustedNetWeight ?? calculateWeightTicketWeightDifference(weightTicket);
  return 0;
};
export const getAllowableWeight = (weightTicket) => {
  return weightTicket.allowableWeight ?? getWeightTicketNetWeight(weightTicket);
};

export const getTotalNetWeightForWeightTickets = (weightTickets = []) => {
  return weightTickets
    ? weightTickets.reduce((prev, curr) => {
        return prev + getWeightTicketNetWeight(curr);
      }, 0)
    : 0;
};

export const calculatePPMShipmentNetWeight = (shipment) => {
  return getTotalNetWeightForWeightTickets(shipment?.ppmShipment?.weightTickets);
};

export const calculateNonPPMShipmentNetWeight = (shipment) => {
  return returnLowestValue(shipment.primeActualWeight, shipment.reweigh?.weight);
};

export const calculateShipmentNetWeight = (shipment) => {
  if (shipment.ppmShipment) {
    return calculatePPMShipmentNetWeight(shipment);
  }
  return calculateNonPPMShipmentNetWeight(shipment);
};
