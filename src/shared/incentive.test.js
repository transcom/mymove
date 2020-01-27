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
});
