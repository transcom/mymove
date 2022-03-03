import * as formatters from './formatters';

describe('formatters', () => {
  describe('format date for customer app', () => {
    it('should format customer date to DD MMM YYYY', () => {
      expect(formatters.formatCustomerDate('2020-09-27T00:00:00Z')).toBe('27 Sep 2020');
    });
    it('should format signature date to YYYY-MM-DD', () => {
      expect(formatters.formatSignatureDate('2020-09-27T00:00:00Z')).toBe('2020-09-27');
    });
  });
  describe('format order type for customer app', () => {
    it('should format order type to be human readable', () => {
      expect(formatters.formatOrderType('PERMANENT_CHANGE_OF_STATION')).toBe('Permanent change of station');
    });
  });

  describe('formatYesNoInputValue', () => {
    it('returns yes for true', () => {
      expect(formatters.formatYesNoInputValue(true)).toBe('yes');
    });
    it('returns no for false', () => {
      expect(formatters.formatYesNoInputValue(false)).toBe('no');
    });
    it('returns null for anything else', () => {
      expect(formatters.formatYesNoInputValue('true')).toBe(null);
      expect(formatters.formatYesNoInputValue('false')).toBe(null);
      expect(formatters.formatYesNoInputValue('')).toBe(null);
      expect(formatters.formatYesNoInputValue({})).toBe(null);
      expect(formatters.formatYesNoInputValue(0)).toBe(null);
      expect(formatters.formatYesNoInputValue(undefined)).toBe(null);
    });
  });

  describe('formatYesNoAPIValue', () => {
    it('returns true for yes', () => {
      expect(formatters.formatYesNoAPIValue('yes')).toBe(true);
    });
    it('returns false for no', () => {
      expect(formatters.formatYesNoAPIValue('no')).toBe(false);
    });
    it('returns undefined for anything else', () => {
      expect(formatters.formatYesNoAPIValue('true')).toBe(undefined);
      expect(formatters.formatYesNoAPIValue('false')).toBe(undefined);
      expect(formatters.formatYesNoAPIValue(true)).toBe(undefined);
      expect(formatters.formatYesNoAPIValue(false)).toBe(undefined);
      expect(formatters.formatYesNoAPIValue('')).toBe(undefined);
      expect(formatters.formatYesNoAPIValue({})).toBe(undefined);
      expect(formatters.formatYesNoAPIValue(0)).toBe(undefined);
      expect(formatters.formatYesNoAPIValue(null)).toBe(undefined);
    });
  });

  describe('formatWeightCWTFromLbs', () => {
    it('returns expected value', () => {
      expect(formatters.formatWeightCWTFromLbs('8000')).toBe('80 cwt');
    });
  });

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

  describe('formatDollarFromMillicents', () => {
    it('returns expected value', () => {
      expect(formatters.formatDollarFromMillicents('80000')).toBe('$0.80');
    });
  });

  describe('formatCents', () => {
    it('formats cents value into local string to 2 decimal places', () => {
      expect(formatters.formatCents(120034)).toEqual('1,200.34');
    });

    it('formats without decimal place when fraction digits are zero', () => {
      expect(formatters.formatCents(120034, 0, 0)).toEqual('1,200');
    });
  });

  describe('formatCentsTruncateWhole', () => {
    it('formats cents value into local string and truncates decimal', () => {
      expect(formatters.formatCentsTruncateWhole(120034)).toEqual('1,200');
    });
  });

  describe('formatDaysInTransit', () => {
    it('returns 0 days when value is null', () => {
      expect(formatters.formatDaysInTransit()).toEqual('0 days');
    });

    it('returns 0 days when value is zero', () => {
      expect(formatters.formatDaysInTransit(0)).toEqual('0 days');
    });

    it('returns 1 day when value is one', () => {
      expect(formatters.formatDaysInTransit(1)).toEqual('1 day');
    });

    it('returns plural when greater than 1', () => {
      expect(formatters.formatDaysInTransit(2)).toEqual('2 days');
    });
  });

  describe('formatDelimitedNumber', () => {
    it('works for simple string numbers', () => {
      expect(formatters.formatDelimitedNumber('500')).toEqual(500);
      expect(formatters.formatDelimitedNumber('1,234')).toEqual(1234);
      expect(formatters.formatDelimitedNumber('12,345,678,901')).toEqual(12345678901);
    });

    it('works for actual numbers', () => {
      expect(formatters.formatDelimitedNumber(500)).toEqual(500);
      expect(formatters.formatDelimitedNumber(1234)).toEqual(1234);
    });

    it('works for non-integers', () => {
      expect(formatters.formatDelimitedNumber('1,234.56')).toEqual(1234.56);
    });
  });
});
