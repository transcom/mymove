import queryString from 'query-string';

import { milmoveLog, MILMOVE_LOG_LEVEL } from 'utils/milmoveLog';

// Simple feature toggling for client-side code.
//
// Environments:
// make client_run -> development
// make client_test -> test
// make client_build -> branches based on hostname, see switch statement below

const defaultFlags = {
  ppm: true,
  documentViewer: true,
  moveInfoComboButton: true,
  sitPanel: true,
  ppmPaymentRequest: true,
  allOrdersTypes: false,
  hhgFlow: false,
  ghcFlow: false,
  markerIO: false,
};

const environmentFlags = {
  development: {
    ...defaultFlags,
    allOrdersTypes: true,
    hhgFlow: true,
    ghcFlow: true,
  },

  test: {
    ...defaultFlags,
  },

  experimental: {
    ...defaultFlags,
    allOrdersTypes: true,
    hhgFlow: true,
    ghcFlow: true,
  },

  staging: {
    ...defaultFlags,
    allOrdersTypes: true,
    hhgFlow: true,
    ghcFlow: true,
    markerIO: true,
  },

  demo: {
    ...defaultFlags,
    allOrdersTypes: true,
    hhgFlow: true,
    ghcFlow: true,
    markerIO: true,
  },

  production: {
    ...defaultFlags,
    sitPanel: false,
  },
};

const validateFlag = (name) => {
  // Warn if the value being fetched was never set.
  if (Object.keys(defaultFlags).indexOf(name) === -1) {
    milmoveLog(MILMOVE_LOG_LEVEL.WARN, `'${name}' is not a valid flag name.`);
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

export const createModifiedSchemaForOrdersTypesFlag = (schema) => {
  return {
    ...schema,
    properties: {
      ...schema.properties,
      orders_type: {
        ...schema.properties.orders_type,
        enum: [schema.properties.orders_type.enum[0]],
      },
    },
  };
};
