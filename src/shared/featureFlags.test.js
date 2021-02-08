import { detectEnvironment, detectFlags } from './featureFlags';

describe('feature flags', () => {
  it('detects the environment correctly', () => {
    expect(detectEnvironment('development', '')).toEqual('development');

    expect(detectEnvironment('test', '')).toEqual('test');

    expect(detectEnvironment('production', 'office.move.mil')).toEqual('production');
    expect(detectEnvironment('production', 'my.move.mil')).toEqual('production');

    expect(detectEnvironment('production', 'office.stg.move.mil')).toEqual('staging');
    expect(detectEnvironment('production', 'my.stg.move.mil')).toEqual('staging');

    expect(detectEnvironment('production', 'office.exp.move.mil')).toEqual('experimental');
    expect(detectEnvironment('production', 'my.exp.move.mil')).toEqual('experimental');

    expect(detectEnvironment('production', 'localhost')).toEqual('development');
    expect(detectEnvironment('production', 'milmovelocal')).toEqual('development');
    expect(detectEnvironment('production', 'officelocal')).toEqual('development');
    expect(detectEnvironment('production', '')).toEqual('development');
  });

  it('merged query string flags into those from the environment', () => {
    const flags = detectFlags('development', 'hostname', '?flag:ppm=false');
    expect(flags.ppm).toEqual(false);
    expect(flags.doesntexist).toEqual(undefined);
  });
});
