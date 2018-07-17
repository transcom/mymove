import { detectEnvironment, detectFlags } from './featureFlags';

describe('feature flags', () => {
  it('detects the environment correctly', () => {
    expect(detectEnvironment('development', '')).toEqual('development');

    expect(detectEnvironment('test', '')).toEqual('test');

    expect(detectEnvironment('production', 'tsp.move.mil')).toEqual(
      'production',
    );

    expect(detectEnvironment('production', 'office.move.mil')).toEqual(
      'production',
    );
    expect(detectEnvironment('production', 'my.move.mil')).toEqual(
      'production',
    );

    expect(detectEnvironment('production', 'tsp.staging.move.mil')).toEqual(
      'staging',
    );
    expect(detectEnvironment('production', 'office.staging.move.mil')).toEqual(
      'staging',
    );
    expect(detectEnvironment('production', 'my.staging.move.mil')).toEqual(
      'staging',
    );

    expect(
      detectEnvironment('production', 'tsp.experimental.move.mil'),
    ).toEqual('experimental');
    expect(
      detectEnvironment('production', 'office.experimental.move.mil'),
    ).toEqual('experimental');
    expect(detectEnvironment('production', 'my.experimental.move.mil')).toEqual(
      'experimental',
    );

    expect(detectEnvironment('production', 'localhost')).toEqual('development');
    expect(detectEnvironment('production', '')).toEqual('development');
  });

  it('merged query string flags into those from the environment', () => {
    const flags = detectFlags('development', 'hostname', '?flag:ppm=false');
    expect(flags.ppm).toEqual(false);
    expect(flags.doesntexist).toEqual(undefined);
  });
});
