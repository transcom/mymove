import * as formatters from './formatters';

describe('formatters', () => {
  describe('formatWeight', () => {
    describe('when formatting a integer weight', () => {
      const weight = 4000;
      const formattedWeight = formatters.formatWeight(weight);
      it('should be be formatted as expected', () => {
        expect(formattedWeight).toEqual('4,000 lbs');
      });
    });
    describe('when formatting a integer weight', () => {
      const weight = '';
      const formattedWeight = formatters.formatWeight(weight);
      it('should be be formatted as expected', () => {
        expect(formattedWeight).toEqual('0 lbs');
      });
    });
  });
  describe('truncateNumber', () => {
    it('should truncate number based on passed variable returning a number string', () => {
      expect(formatters.truncateNumber(500000)).toEqual('50');
      expect(formatters.truncateNumber(500000, 0)).toEqual('50');
      expect(formatters.truncateNumber(505555, 1)).toEqual('50.5');
      expect(formatters.truncateNumber(500000, 2)).toEqual('50.00');
      expect(formatters.truncateNumber(505555, 2)).toEqual('50.55');
      expect(formatters.truncateNumber(505555, 4)).toEqual('50.5555');
      expect(formatters.truncateNumber(505555, 10)).toEqual('50.5555000000');
      expect(formatters.truncateNumber(500, 2)).toEqual('0.05');
      expect(formatters.truncateNumber(500, 0)).toEqual('0');
      expect(formatters.truncateNumber(null)).toEqual(null);
      // negitive numbers act a bit different than expected due to floor()
      expect(formatters.truncateNumber(-505555, 3)).toEqual('-50.556');
    });
  });
  describe('addCommasToNumberString', () => {
    it('should truncate number based on passed variable returning a number string', () => {
      expect(formatters.addCommasToNumberString(5000)).toEqual('5,000');
      expect(formatters.addCommasToNumberString('500000000.0001')).toEqual('500,000,000.0001');
      expect(formatters.addCommasToNumberString('500000000')).toEqual('500,000,000');
      expect(formatters.addCommasToNumberString('5000')).toEqual('5,000');
      expect(formatters.addCommasToNumberString('-5000')).toEqual('-5,000');
      expect(formatters.addCommasToNumberString('500')).toEqual('500');
      expect(formatters.addCommasToNumberString('0')).toEqual('0');
      expect(formatters.addCommasToNumberString('0', 2)).toEqual('0.00');
    });
  });
});
