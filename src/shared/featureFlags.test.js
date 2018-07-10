import { detectEnvironment, feature, override, reset } from './featureFlags';

describe('feature flags', () => {
  it('detects the environment correctly', () => {
    expect(detectEnvironment()).toEqual('test');
  });

  it('returns the correct values for the test environment', () => {
    expect(feature('ppm')).toEqual(true);
    expect(feature('justForTesting')).toEqual(false);
    expect(feature('doesntexist')).toEqual(undefined);
  });

  it('can override flags using override', () => {
    override('ppm', false);
    expect(feature('ppm')).toEqual(false);
  });

  it('can reset after overriding a flag', () => {
    override('ppm', false);
    reset();
    expect(feature('ppm')).toEqual(true);
  });
});
