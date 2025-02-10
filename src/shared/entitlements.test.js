import { selectEntitlements } from './entitlements';

describe('entitlements', () => {
  describe('when I am not logged in', () => {
    it('should return an empty object', () => {
      const entitlements = selectEntitlements();
      expect(entitlements).toEqual({});
    });
  });
  describe('when I have dependents', () => {
    describe('when my spouse has pro gear', () => {
      it('should include spouse progear', () => {
        const entitlements = selectEntitlements(
          {
            authorizedWeight: 8000,
            entitlement: {
              proGear: 2000,
              proGearSpouse: 500,
            },
          },
          true,
          true,
        );
        expect(entitlements).toEqual({
          proGear: 2000,
          proGearSpouse: 500,
          sum: 10500,
          weight: 8000,
        });
      });
    });
    describe('when my spouse does not have pro gear', () => {
      it('should not include spouse progear', () => {
        const entitlements = selectEntitlements(
          {
            authorizedWeight: 8000,
            entitlement: {
              proGear: 2000,
              proGearSpouse: 0,
            },
          },
          true,
          false,
        );
        expect(entitlements).toEqual({
          proGear: 2000,
          proGearSpouse: 0,
          sum: 10000,
          weight: 8000,
        });
      });
    });
  });
  describe("when I don't have dependents", () => {
    it('should exclude spouse progear', () => {
      const entitlements = selectEntitlements({
        authorizedWeight: 5000,
        entitlement: {
          proGear: 2000,
          proGearSpouse: 500,
        },
      });
      expect(entitlements).toEqual({
        proGear: 2000,
        proGearSpouse: 0,
        sum: 7000,
        weight: 5000,
      });
    });
  });
});
