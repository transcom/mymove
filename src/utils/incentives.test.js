import {
  hasShortHaulError,
  getIncentiveRange,
  calculateMaxAdvance,
  calculateMaxAdvanceAndFormatAdvanceAndIncentive,
} from './incentives';

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

describe('getIncentiveRange', () => {
  it('should return the formatted range from the PPM if the PPM values exist', () => {
    expect(
      getIncentiveRange(
        {
          incentive_estimate_min: 1000,
          incentive_estimate_max: 2400,
        },
        { range_min: 1400, range_max: 2300 },
      ),
    ).toBe('$10.00 - 24.00');
  });

  it('should return the formatted range from the estimate if the PPM values do not exist', () => {
    expect(getIncentiveRange({}, { range_min: 1400, range_max: 2300 })).toBe('$14.00 - 23.00');
  });

  it('should return an empty string if no values exist', () => {
    expect(getIncentiveRange({}, {})).toBe('');
    expect(
      getIncentiveRange(
        {
          incentive_estimate_max: '',
          incentive_estimate_min: null,
        },
        {
          range_min: 0,
          range_max: undefined,
        },
      ),
    ).toBe('');
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
