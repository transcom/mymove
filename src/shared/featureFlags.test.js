import { feature, flags } from './featureFlags';

describe('feature flags', () => {
  it('returns the correct values for the test environment', () => {
    expect(feature('ppm')).toEqual(true);
    expect(feature('justForTesting')).toEqual(false);
    expect(feature('doesntexist')).toEqual(undefined);
  });
});
