import queryString from 'query-string';
import { forEach } from 'lodash';

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

const environmentFlags = {
  development: Object.assign({}, defaultFlags),

  test: Object.assign({}, defaultFlags),

  staging: Object.assign({}, defaultFlags, {
    hhg: false,
  }),

  production: Object.assign({}, defaultFlags, {
    hhg: false,
  }),
};

export function flagsFromURL(search) {
  const params = queryString.parse(search);
  let flags = {};

  forEach(params, function(value, key) {
    let [prefix, name] = key.split(':');
    if (prefix === 'flag' && name.length > 0) {
      if (validateFlag(name)) {
        // name is validated by the previous line
        // eslint-disable-next-line security/detect-object-injection
        flags[name] = value === 'true';
      }
    }
  });
  return flags;
}

// Return the name of the current envirnonment as a string.
export function detectEnvironment(nodeEnv, host) {
  if (nodeEnv !== 'production') {
    return nodeEnv;
  }

  // If we've built the app, then use the hostname to determine what the
  // environment is.
  const domain = host;
  switch (domain) {
    case 'office.move.mil':
    case 'my.move.mil':
      return 'production';
    case 'office-staging.move.mil':
    case 'my-staging.move.mil':
      return 'staging';
    // TODO add experimental
    default:
      return 'development';
  }
}

function validateFlag(name) {
  // Warn if the value is undefined, indicating that a value is being fetched
  // that was never set.
  if (Object.keys(defaultFlags).indexOf(name) === -1) {
    if (console && console.warn) {
      console.warn(`'${name}' is not a valid flag name.`);
    }
    return false;
  }
  return true;
}

export function detectFlags(nodeEnv, host, search) {
  let env = detectEnvironment(nodeEnv, host);
  // env can only be one of the values hard-coded into detectEnvironment()
  // eslint-disable-next-line security/detect-object-injection
  return Object.assign({}, environmentFlags[env], flagsFromURL(search));
}
