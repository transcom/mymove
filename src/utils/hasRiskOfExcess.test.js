import hasRiskOfExcess from './hasRiskOfExcess';

describe('hasRiskOfExcess', () => {
  it('returns true when estimated weight is 90% of weight allowancew', () => {
    expect(hasRiskOfExcess(90, 100)).toEqual(true);
  });

  it('returns true when estimated weight is more than 90% of weight allowancew', () => {
    expect(hasRiskOfExcess(91, 100)).toEqual(true);
  });

  it('returns false when estimated weight is less than 90% of weight allowancew', () => {
    expect(hasRiskOfExcess(89, 100)).toEqual(false);
  });

  it('returns false when estimated weight is undefined', () => {
    expect(hasRiskOfExcess(undefined, 100)).toEqual(false);
  });

  it('returns false when estimated weight is zero', () => {
    expect(hasRiskOfExcess(0, 100)).toEqual(false);
  });
});
