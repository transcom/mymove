import { numOfDaysBetweenDates, selectDateFieldByStatus, selectDatePrefixByStatus } from './dates';

describe('numOfDaysBetweenDates', () => {
  it('should return 5 for number of days between Aug 5th and Aug 10', () => {
    expect(numOfDaysBetweenDates('2022-08-01', '2022-08-06')).toBe(5);
  });
  it('should return 6 for number of days between Aug 31 and Sept 6th', () => {
    expect(numOfDaysBetweenDates('2022-08-31', '2022-09-06')).toBe(6);
  });
});

describe('selectDatePrefixByStatus', () => {
  it('should return "Date requested" for a SUBMITTED status', () => {
    expect(selectDatePrefixByStatus('SUBMITTED')).toEqual('Date requested');
  });
  it('should return "Date approved" for a APPROVED status', () => {
    expect(selectDatePrefixByStatus('APPROVED')).toEqual('Date approved');
  });
  it('should return "Date rejected" for a REJECTED status', () => {
    expect(selectDatePrefixByStatus('REJECTED')).toEqual('Date rejected');
  });
  it('should return "Date approved" for a Move Task Order Approved status', () => {
    expect(selectDatePrefixByStatus('Move Task Order Approved')).toEqual('Date approved');
  });
  it('should return "Date rejected" for a Move Task Order Rejected status', () => {
    expect(selectDatePrefixByStatus('Move Task Order Rejected')).toEqual('Date rejected');
  });
  it('should return "Date requested" as default', () => {
    expect(selectDatePrefixByStatus('noMatch')).toEqual('Date requested');
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
  it('should return approvedAt for a Move Task Order Approved status', () => {
    expect(selectDateFieldByStatus('Move Task Order Approved')).toEqual('approvedAt');
  });
  it('should return rejectedAt for a Move Task Order Rejected status', () => {
    expect(selectDateFieldByStatus('Move Task Order Rejected')).toEqual('rejectedAt');
  });
  it('should return createdAt as default', () => {
    expect(selectDateFieldByStatus('noMatch')).toEqual('createdAt');
  });
});
