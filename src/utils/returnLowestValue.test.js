import returnLowestValue from './returnLowestValue';

describe('returnLowestValue', () => {
  it('returns lower value of two numbers', () => {
    expect(returnLowestValue(99, 100)).toEqual(99);
  });

  it('returns number if only one valid value is passed', () => {
    expect(returnLowestValue(100, null)).toEqual(100);
    expect(returnLowestValue(null, 100)).toEqual(100);
  });

  it('returns null if two falsy values are passed', () => {
    expect(returnLowestValue(null, null)).toEqual(null);
  });
});
