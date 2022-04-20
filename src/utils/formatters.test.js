import * as formatters from './formatters';

import PAYMENT_REQUEST_STATUS from 'constants/paymentRequestStatus';

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

  describe('formatLabelReportByDate', () => {
    it('returns the correct label for RETIREMENT', () => {
      expect(formatters.formatLabelReportByDate('RETIREMENT')).toEqual('Date of retirement');
    });
    it('returns the correct label for SEPARATION', () => {
      expect(formatters.formatLabelReportByDate('SEPARATION')).toEqual('Date of separation');
    });
    it('returns a default label for all other values', () => {
      expect(formatters.formatLabelReportByDate('test')).toEqual('Report by date');
    });
  });

  describe('toDollarString', () => {
    it('returns string representation of a dollar', () => {
      expect(formatters.toDollarString(1234.12)).toEqual('$1,234.12');
    });
  });

  describe('formatDate', () => {
    it('should be formatted as expected', () => {
      const inputFormat = 'MMM-DD-YY';
      const formattedDate = formatters.formatDate('Nov-11-19', inputFormat, 'DD-MMM-YY', 'en', true);
      expect(formattedDate).toBe('11-Nov-19');
    });

    it('should be invalid with unexpected input and strict mode on', () => {
      const inputFormat = 'MMM-DD-YY';
      const formattedDate = formatters.formatDate('Nov-11-1999', inputFormat, 'DD-MMM-YY', 'en', true);
      expect(formattedDate).toBe('Invalid date');
    });

    it('should default to DD-MMM-YY ouptut format', () => {
      const inputFormat = 'MMM-DD-YY';
      expect(formatters.formatDate('Nov-11-99', inputFormat)).toBe('11-Nov-99');
    });
  });

  describe('formatDateFromIso', () => {
    it('should be formatted as expected', () => {
      expect(formatters.formatDateFromIso('2020-08-11T21:00:59.126987Z', 'DD MMM YYYY')).toBe('11 Aug 2020');
    });
  });

  describe('formatToOrdinal', () => {
    it('returns the ordinal corresponding to an int', () => {
      expect(formatters.formatToOrdinal(1)).toEqual('1st');
      expect(formatters.formatToOrdinal(2)).toEqual('2nd');
      expect(formatters.formatToOrdinal(3)).toEqual('3rd');
      expect(formatters.formatToOrdinal(4)).toEqual('4th');
    });
  });

  describe('paymentRequestStatusReadable', () => {
    it('returns expected string for PENDING', () => {
      expect(formatters.paymentRequestStatusReadable(PAYMENT_REQUEST_STATUS.PENDING)).toEqual('Payment requested');
    });

    it('returns expected string for REVIEWED', () => {
      expect(formatters.paymentRequestStatusReadable(PAYMENT_REQUEST_STATUS.REVIEWED)).toEqual('Reviewed');
    });

    it('returns expected string for SENT_TO_GEX', () => {
      expect(formatters.paymentRequestStatusReadable(PAYMENT_REQUEST_STATUS.SENT_TO_GEX)).toEqual('Reviewed');
    });

    it('returns expected string for RECEIVED_BY_GEX', () => {
      expect(formatters.paymentRequestStatusReadable(PAYMENT_REQUEST_STATUS.RECEIVED_BY_GEX)).toEqual('Reviewed');
    });

    it('returns expected string for PAID', () => {
      expect(formatters.paymentRequestStatusReadable(PAYMENT_REQUEST_STATUS.PAID)).toEqual('Paid');
    });

    it('returns expected string for EDI_ERROR', () => {
      expect(formatters.paymentRequestStatusReadable(PAYMENT_REQUEST_STATUS.EDI_ERROR)).toEqual('EDI error');
    });

    it('returns expected string for DEPRECATED', () => {
      expect(formatters.paymentRequestStatusReadable(PAYMENT_REQUEST_STATUS.DEPRECATED)).toEqual('Deprecated');
    });
  });

  describe('filenameFromPath', () => {
    it('returns last portion of path with default delimiter', () => {
      expect(formatters.filenameFromPath('/home/user/folder/.hidden/My Long Filename.sql')).toEqual(
        'My Long Filename.sql',
      );
    });

    it('returns original filename if no path is included', () => {
      expect(formatters.filenameFromPath('Just-A-gnarly_filemame(0) DRAFT.v2.docx')).toEqual(
        'Just-A-gnarly_filemame(0) DRAFT.v2.docx',
      );
    });
  });

  describe('formatAgeToDays', () => {
    it('returns expected string less than 1 day', () => {
      expect(formatters.formatAgeToDays(0.99)).toEqual('Less than 1 day');
    });

    it('returns expected string for 1 day', () => {
      expect(formatters.formatAgeToDays(1.5)).toEqual('1 day');
    });

    it('returns expected string greater than 1 day', () => {
      expect(formatters.formatAgeToDays(2.99)).toEqual('2 days');
    });
  });
});
