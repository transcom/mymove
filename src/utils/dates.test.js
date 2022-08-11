import numOfDaysBetweenDates from './dates';

describe('numOfDaysBetweenDates', () => {
  it('should return 5 for number of days between Aug 5th and Aug 10', () => {
    expect(numOfDaysBetweenDates('2022-08-01', '2022-08-06')).toBe(5);
  });
  it('should return 6 for number of days between Aug 31 and Sept 6th', () => {
    expect(numOfDaysBetweenDates('2022-08-31', '2022-09-06')).toBe(6);
  });
});
