import {
  shipmentIsOverweight,
  calcWeightRequested,
  calcTotalBillableWeight,
  calcTotalEstimatedWeight,
} from './shipmentWeights';

describe('shipmentWeights utils', () => {
  describe('shipmentIsOverweight', () => {
    it('returns true when the shipment weight is over 110% of the estimated weight', () => {
      expect(shipmentIsOverweight(100, 111)).toEqual(true);
    });

    it('returns false when shipment weight is less than  110% of the estimated weight', () => {
      expect(shipmentIsOverweight(100, 101)).toEqual(false);
    });

    it('returns false when estimated weight is undefined', () => {
      expect(shipmentIsOverweight(undefined, 100)).toEqual(false);
    });
  });
  describe('calcWeightRequested', () => {
    it('returns sum of actual weights if no reweigh weights', () => {
      const mtoShipments = [
        { billableWeightCap: 1000, primeActualWeight: 300 },
        { billableWeightCap: 2000, primeActualWeight: 400 },
        { billableWeightCap: 3000, primeActualWeight: 300 },
      ];
      expect(calcWeightRequested(mtoShipments)).toEqual(1000);
    });

    it('return sum of smaller value between reweigh weights and actual weight', () => {
      const mtoShipments = [
        { billableWeightCap: 1000, primeActualWeight: 300, reweigh: { weight: 100 } },
        { billableWeightCap: 2000, primeActualWeight: 400, reweigh: { weight: 1000 } },
        { billableWeightCap: 3000, primeActualWeight: 300, reweigh: { weight: 200 } },
      ];
      expect(calcWeightRequested(mtoShipments)).toEqual(700);
    });
  });
  describe('calcTotalBillableWeight', () => {
    it('returns sum of billable weight if provided', () => {
      const mtoShipments = [
        { billableWeightCap: 1000, primeActualWeight: 300, reweigh: { weight: 100 } },
        { billableWeightCap: 2000, primeActualWeight: 400, reweigh: { weight: 1000 } },
        { billableWeightCap: 3000, primeActualWeight: 300, reweigh: { weight: 200 } },
      ];
      expect(calcTotalBillableWeight(mtoShipments)).toEqual(6000);
    });
    it('returns sum of actual weights if there are no billable weights and reweighs', () => {
      const mtoShipments = [{ primeActualWeight: 300 }, { primeActualWeight: 400 }, { primeActualWeight: 300 }];
      expect(calcTotalBillableWeight(mtoShipments)).toEqual(1000);
    });

    it('returns the sum of smaller value between reweigh weights and actual weight', () => {
      const mtoShipments = [
        { primeActualWeight: 300, reweigh: { weight: 100 } },
        { primeActualWeight: 400, reweigh: { weight: 1000 } },
        { primeActualWeight: 300, reweigh: { weight: 200 } },
      ];
      expect(calcTotalBillableWeight(mtoShipments)).toEqual(700);
    });
  });
  describe('calcTotalEstimatedWeight', () => {
    it('returns the sum of shipments estimated weight', () => {
      const mtoShipments = [
        { primeEstimatedWeight: 1000, primeActualWeight: 300, reweigh: { weight: 100 } },
        { primeEstimatedWeight: 2000, primeActualWeight: 400, reweigh: { weight: 1000 } },
        { primeEstimatedWeight: 7000, primeActualWeight: 300, reweigh: { weight: 200 } },
      ];
      expect(calcTotalEstimatedWeight(mtoShipments)).toEqual(10000);
    });
  });
});
