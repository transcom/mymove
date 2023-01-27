import { parseDate, formatDate, formatDateForSwagger, formatDateTime } from './dates';
import moment from 'moment';
describe('dates', () => {
  describe('parseDate', () => {
    describe('when parsing a date that does not match the allowed date formats', () => {
      const result = parseDate('8');
      it('should return undefined', () => {
        expect(result).toBeUndefined();
      });
    });
    describe('when parsing a date that does match the allowed date formats', () => {
      const result = parseDate('8-23-2019');
      it('should return a Date that matches that string', () => {
        expect(result).toEqual(new moment('2019-08-23T00:00:00.000Z').toDate());
      });
    });
  });
  describe('formatDate', () => {
    describe('when formatting a date that does not match the allowed date formats', () => {
      it('should return "invalid date"', () => {
        const result = formatDate('8');
        expect(result).toEqual('Invalid date');
      });
      it('should return "invalid date"', () => {
        const result = formatDate();
        expect(result).toEqual('Invalid date');
      });
    });
    describe('when parsing a date that does the match allowed date formats', () => {
      const result = formatDate('8-23-2019');
      it('should return 8/23/2019', () => {
        expect(result).toEqual('8/23/2019');
      });
    });
  });
  describe('formatDateForSwagger', () => {
    describe('when formatting a date that does not match the allowed date formats', () => {
      it('should return invalid date', () => {
        const result = formatDateForSwagger('8');
        expect(result).toEqual('Invalid date');
      });
      it('should return undefined', () => {
        const result = formatDateForSwagger();
        expect(result).toEqual(undefined);
      });
    });
    describe('when parsing a date that does the match allowed date formats', () => {
      const result = formatDateForSwagger('8-23-2019');
      it('should return a date in the format swagger accepts', () => {
        expect(result).toEqual('2019-08-23');
      });
    });
  });
  describe('formatDateTime', () => {
    describe('when parsing a date that does the match allowed date formats', () => {
      const result = formatDateTime('8-23-2019');
      it('should return a date in the format swagger accepts', () => {
        expect(result).toEqual('2019-08-23T00:00:00+00:00');
      });
    });
  });
});
