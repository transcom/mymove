import * as formatters from './formatters';
import moment from 'moment';

describe('formatters', () => {
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
