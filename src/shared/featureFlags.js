// Simple feature toggling for client-side code.
//
// Environments:
// make client_run -> development
// make client_test -> test
// make client_build -> branches based on hostname, see switch statement below

const defaultFlags = {
  ppm: true,
  hhg: true,
};

export const flags = {
  development: Object.assign({}, defaultFlags),
  test: Object.assign({}, defaultFlags, {
    justForTesting: false,
  }),
  staging: Object.assign({}, defaultFlags),
  experimental: Object.assign({}, defaultFlags),

  production: Object.assign({}, defaultFlags, {
    hhg: false,
  }),
};

// Return the name of the current envirnonment as a string.
function detectEnvironment() {
  const nodeEnv = process.env['NODE_ENV'];

  if (nodeEnv !== 'production') {
    return nodeEnv;
  }

  // If we've built the app, then use the hostname to determine what the
  // environment is.
  const domain = window.location.host;
  switch (domain) {
    case 'office.move.mil':
    case 'my.move.mil':
      return 'production';
      break;
    case 'office-staging.move.mil':
    case 'my-staging.move.mil':
      return 'staging';
      return;
    // TODO add experimental
    default:
      return 'development';
  }
}

// Return true or false based on if the requested feature is enabled in this
// environment.
export function feature(name) {
  let env = detectEnvironment();
  let value = flags[env][name];

  // Warn if the value is undefined, indicating that a value is being fetched
  // that was never set.
  if (
    typeof value === 'undefined' &&
    env !== 'test' &&
    console &&
    console.warn
  ) {
    console.warn(
      `Value for flag '${name}' in environment '${env}' is undefined.`,
    );
  }
  return value;
}
