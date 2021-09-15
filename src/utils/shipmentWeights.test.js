import { shipmentIsOverweight } from './shipmentWeights';

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
});
