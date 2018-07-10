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

  experimental: Object.assign({}, defaultFlags),

  staging: Object.assign({}, defaultFlags, {
    hhg: false,
  }),

  production: Object.assign({}, defaultFlags, {
    hhg: false,
  }),
};

let overrides = {};

// Return the name of the current envirnonment as a string.
export function detectEnvironment() {
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

export function override(name, value) {
  if (Object.keys(defaultFlags).indexOf(name) === -1) {
    throw `Flag '${name}' is not defined. Add it to defaultFlags if you wish to use it.`;
  }
  overrides[name] = value;
}

export function reset() {
  overrides = {};
}

// Return true or false based on if the requested feature is enabled in this
// environment.
export function feature(name) {
  if (Object.keys(overrides).indexOf(name) !== -1) {
    return overrides[name];
  }

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
