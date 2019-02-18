import * as preApprovals from './preApprovals';

describe('preApprovals', () => {
  describe('isNewAccessorial 105B or 105E', () => {
    it('should return true if new accessorial, false if old accessorial', () => {
      const item105BOld = { tariff400ng_item: { code: '105B' } };
      const item105BNew = { tariff400ng_item: { code: '105B' }, crate_dimensions: { length: 1, height: 1, width: 1 } };
      const item105EOld = { tariff400ng_item: { code: '105E' } };
      const item105ENew = { tariff400ng_item: { code: '105E' }, crate_dimensions: { length: 1, height: 1, width: 1 } };

      const itemNull = null;

      expect(preApprovals.isNewAccessorial(item105BOld)).toEqual(false);
      expect(preApprovals.isNewAccessorial(item105BNew)).toEqual(true);
      expect(preApprovals.isNewAccessorial(item105EOld)).toEqual(false);
      expect(preApprovals.isNewAccessorial(item105ENew)).toEqual(true);
      expect(preApprovals.isNewAccessorial(itemNull)).toEqual(false);
    });
  });
});
