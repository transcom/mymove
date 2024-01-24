import { hasShortHaulError, calculateMaxAdvance, calculateMaxAdvanceAndFormatAdvanceAndIncentive } from './incentives';

describe('hasShortHaulError', () => {
  it('should return true for 409 - move under 50 miles', () => {
    expect(hasShortHaulError({ statusCode: 409 })).toBe(true);
  });
  it('should return false for 404 - rate data missing', () => {
    expect(hasShortHaulError({ statusCode: 404 })).toBe(false);
  });
  it('should return false if error undefined', () => {
    expect(hasShortHaulError()).toBe(false);
  });
});

describe('calculateMaxAdvance', () => {
  it.each([
    [100000, 60000],
    [100005, 60003],
    [100100, 60060],
  ])('should return the expected max advance', (incentive, expectedMaxAdvance) => {
    expect(calculateMaxAdvance(incentive)).toBe(expectedMaxAdvance);
  });
});

describe('calculateMaxAdvanceAndFormatAdvanceAndIncentive', () => {
  it.each([
    [100000, 600, '600', '1,000'],
    [100005, 600, '600', '1,000'],
    [100100, 600, '600', '1,001'],
  ])(
    'should return the expected max advance and incentive values - incentive (in cents): %s',
    (incentive, expectedMaxAdvance, expectedFormattedAdvance, expectedFormattedIncentive) => {
      expect(calculateMaxAdvanceAndFormatAdvanceAndIncentive(incentive)).toStrictEqual({
        maxAdvance: expectedMaxAdvance,
        formattedMaxAdvance: expectedFormattedAdvance,
        formattedIncentive: expectedFormattedIncentive,
      });
    },
  );
});
