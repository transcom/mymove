import queryString from 'query-string';

import { getBooleanFeatureFlagUnauthenticated, getBooleanFeatureFlagForUser } from '../services/internalApi';
import { getBooleanFeatureFlagUnauthenticatedOffice } from '../services/ghcApi';

import { milmoveLogger } from 'utils/milmoveLog';

// Simple feature toggling for client-side code.
//
// Environments:
// make client_run -> development
// make client_test -> test
// make client_build -> branches based on hostname, see switch statement below

// Please do not utilize these default / environment flags. These flags have been deprecated in place of Flipt.
// Refer to the feature flag documentation within this project's docs.
const defaultFlags = {
  ppm: true,
  documentViewer: true,
  moveInfoComboButton: true,
  sitPanel: true,
  ppmPaymentRequest: true,
  allOrdersTypes: true,
  hhgFlow: true,
  ghcFlow: true,
  markerIO: false,
};

const environmentFlags = {
  development: {
    ...defaultFlags,
  },

  test: {
    ...defaultFlags,
  },

  experimental: {
    ...defaultFlags,
  },

  staging: {
    ...defaultFlags,
    markerIO: true,
  },

  demo: {
    ...defaultFlags,
    markerIO: true,
  },

  production: {
    ...defaultFlags,
  },

  loadtest: {
    ...defaultFlags,
  },
};

const validateFlag = (name) => {
  // Warn if the value being fetched was never set.
  if (Object.keys(defaultFlags).indexOf(name) === -1) {
    milmoveLogger.warn(`'${name}' is not a valid flag name.`);
    return false;
  }
  return true;
};

export function flagsFromURL(search) {
  const params = queryString.parse(search);
  return Object.entries(params).reduce((mem, pair) => {
    const [key, value] = pair;
    const [prefix, name] = key.split(':');
    if (prefix === 'flag' && name.length > 0 && validateFlag(name)) {
      return {
        ...mem,
        [name]: value === 'true',
      };
    }
    return mem;
  }, {});
}

// Return the name of the current environment as a string.
export function detectEnvironment(nodeEnv, host) {
  if (nodeEnv !== 'production') {
    return nodeEnv;
  }

  // If we've built the app, then use the hostname to determine what the
  // environment is.
  switch (host) {
    case 'office.move.mil':
    case 'my.move.mil':
    case 'admin.move.mil:':
      return 'production';
    case 'office.stg.move.mil':
    case 'my.stg.move.mil':
    case 'admin.stg.move.mil':
      return 'staging';
    case 'my.exp.move.mil':
    case 'office.exp.move.mil':
    case 'admin.exp.move.mil':
      return 'experimental';
    case 'my.demo.dp3.us':
    case 'office.demo.dp3.us':
    case 'admin.demo.dp3.us':
      return 'demo';
    default:
      return 'development';
  }
}

export function detectFlags(nodeEnv, host, search) {
  const env = detectEnvironment(nodeEnv, host);
  // env can only be one of the values hard-coded into detectEnvironment()
  return {
    ...environmentFlags[env],
    ...flagsFromURL(search),
  };
}

// isBooleanFlagEnabled returns the Flipt feature flag value
export async function isBooleanFlagEnabled(flagKey) {
  return getBooleanFeatureFlagForUser(flagKey, {})
    .then((result) => {
      if (result && typeof result.match !== 'undefined') {
        // Found feature flag, "match" is its boolean value
        return result.match;
      }
      throw new Error(`feature flag is undefined ${flagKey}`);
    })
    .catch((error) => {
      // On error, log it and then just return false setting it to be disabled.
      // No need to return it for extra handling.
      milmoveLogger.error(error);
      return false;
    });
}

// isBooleanFlagEnabledUnauthenticated returns the Flipt feature flag value
// only used within the customer app and for unauthenticated users
export async function isBooleanFlagEnabledUnauthenticated(flagKey) {
  return getBooleanFeatureFlagUnauthenticated(flagKey, {})
    .then((result) => {
      if (result && typeof result.match !== 'undefined') {
        return result.match;
      }
      throw new Error(`feature flag is undefined ${flagKey}`);
    })
    .catch((error) => {
      milmoveLogger.error(error);
      return false;
    });
}

// isBooleanFlagEnabledUnauthenticated returns the Flipt feature flag value
// only used within the office app and for unauthenticated users
export async function isBooleanFlagEnabledUnauthenticatedOffice(flagKey) {
  return getBooleanFeatureFlagUnauthenticatedOffice(flagKey, {})
    .then((result) => {
      if (result && typeof result.match !== 'undefined') {
        return result.match;
      }
      throw new Error(`feature flag is undefined ${flagKey}`);
    })
    .catch((error) => {
      milmoveLogger.error(error);
      return false;
    });
}

export function isCounselorMoveCreateEnabled() {
  const flagKey = 'counselor_move_create';
  return getBooleanFeatureFlagForUser(flagKey, {})
    .then((result) => {
      if (result && typeof result.match !== 'undefined') {
        // Found feature flag, "match" is its boolean value
        return result.match;
      }
      throw new Error('counselor move creation feature flag is undefined');
    })
    .catch((error) => {
      // On error, log it and then just return false setting it to be disabled.
      // No need to return it for extra handling.
      milmoveLogger.error(error);
      return false;
    });
}
