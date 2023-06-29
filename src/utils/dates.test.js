import { numOfDaysBetweenDates, selectDateFieldByStatus } from './dates';

describe('numOfDaysBetweenDates', () => {
  it('should return 5 for number of days between Aug 5th and Aug 10', () => {
    expect(numOfDaysBetweenDates('2022-08-01', '2022-08-06')).toBe(5);
  });
  it('should return 6 for number of days between Aug 31 and Sept 6th', () => {
    expect(numOfDaysBetweenDates('2022-08-31', '2022-09-06')).toBe(6);
  });
});

describe('selectDateFieldByStatus', () => {
  it('should return createdAt for a SUBMITTED status', () => {
    expect(selectDateFieldByStatus('SUBMITTED')).toEqual('createdAt');
  });
  it('should return approvedAt for a APPROVED status', () => {
    expect(selectDateFieldByStatus('APPROVED')).toEqual('approvedAt');
  });
  it('should return rejectedAt for a REJECTED status', () => {
    expect(selectDateFieldByStatus('REJECTED')).toEqual('rejectedAt');
  });
  it('should return createdAt as default', () => {
    expect(selectDateFieldByStatus('noMatch')).toEqual('createdAt');
  });
});
