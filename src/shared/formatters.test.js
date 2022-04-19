import * as formatters from './formatters';
import moment from 'moment';

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
