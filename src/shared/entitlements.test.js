import { getEntitlements } from './entitlements';
describe('entitlements', () => {
  describe('when I have dependents', () => {
    it('should include spouse progear', () => {
      const entitlements = getEntitlements(`E_2`, true);
      expect(entitlements).toEqual({
        pro_gear: 2000,
        pro_gear_spouse: 500,
        sum: 10500,
        weight: 8000,
      });
    });
  });
  describe("when I don't have dependents", () => {
    it('should exclude spouse progear', () => {
      const entitlements = getEntitlements(`E_2`);
      expect(entitlements).toEqual({
        pro_gear: 2000,
        pro_gear_spouse: 0,
        sum: 7000,
        weight: 5000,
      });
    });
  });
});
