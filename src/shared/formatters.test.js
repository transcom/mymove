import * as formatters from './formatters';
import moment from 'moment';

import PAYMENT_REQUEST_STATUS from 'constants/paymentRequestStatus';

describe('formatters', () => {
  describe('truncateNumber', () => {
    it('should truncate number based on passed variable returning a number string', () => {
      expect(formatters.truncateNumber(50)).toEqual('50');
      expect(formatters.truncateNumber(50.0, 0)).toEqual('50');
      expect(formatters.truncateNumber(50.5555, 1)).toEqual('50.5');
      expect(formatters.truncateNumber(50, 2)).toEqual('50.00');
      expect(formatters.truncateNumber(50.5555, 2)).toEqual('50.55');
      expect(formatters.truncateNumber(50.5555, 4)).toEqual('50.5555');
      expect(formatters.truncateNumber(50.5555, 10)).toEqual('50.5555000000');
      expect(formatters.truncateNumber(0.05, 2)).toEqual('0.05');
      expect(formatters.truncateNumber(0.05, 0)).toEqual('0');
      expect(formatters.truncateNumber(null)).toEqual(null);
      expect(formatters.truncateNumber('50.55', 1)).toEqual('50.5');
      // negitive numbers act a bit different than expected due to floor()
      expect(formatters.truncateNumber(-50.5555, 3)).toEqual('-50.556');
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

  describe('formatDateTimeWithTZ', () => {
    it('should include the timezone shortcode', () => {
      const formattedDate = formatters.formatDateTimeWithTZ(new Date());
      expect(formattedDate).toMatch(/\d{2}-\w{3}-\d{2} \d{2}:\d{2} \w{2,3}/);
    });
  });

  describe('formatTimeAgo', () => {
    it('should account for 1 minute correctly', () => {
      let time = new Date();
      let formattedTime = formatters.formatTimeAgo(time);

      expect(formattedTime).toEqual('a few seconds ago');

      time = moment().subtract(1, 'minute')._d;
      formattedTime = formatters.formatTimeAgo(time);

      expect(formattedTime).toEqual('1 min ago');
    });
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
