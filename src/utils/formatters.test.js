import * as formatters from './formatters';

describe('formatters', () => {
  describe('format date for customer app', () => {
    it('should format customer date to DD MMM YYYY', () => {
      expect(formatters.formatCustomerDate('Sep-27-20')).toBe('27 Sep 2020');
    });
  });
  describe('format order type for customer app', () => {
    it('should format order type to be human readable', () => {
      expect(formatters.formatOrderType('PERMANENT_CHANGE_OF_STATION')).toBe('Permanent change of station');
    });
  });
});
