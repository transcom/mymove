import * as incentive from './incentive';

describe('incentive', () => {
  describe('Check for short haul error', () => {
    it('should return true for 409 - move under 50 miles', () => {
      expect(incentive.hasShortHaulError({ statusCode: 409 })).toBe(true);
    });
    it('should return false for 404 - rate data missing', () => {
      expect(incentive.hasShortHaulError({ statusCode: 404 })).toBe(false);
    });
    it('should return false if error undefined', () => {
      expect(incentive.hasShortHaulError()).toBe(false);
    });
  });
  describe('Check format for incentive range', () => {
    it('should reutrn range', () => {
      expect(incentive.formatIncentiveRange({ incentive_estimate_min: 1000, incentive_estimate_max: 2000 })).toEqual(
        '$10.00 - 20.00',
      );
      expect(
        incentive.formatIncentiveRange({
          currentPpm: { incentive_estimate_min: 30000, incentive_estimate_max: 40000 },
        }),
      ).toEqual('$300.00 - 400.00');
    });
  });
});
