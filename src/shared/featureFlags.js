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
  documentViewer: true,
  moveInfoComboButton: true,
  sitPanel: true,
  ppmPaymentRequest: true,
  allOrdersTypes: false,
  hhgFlow: false,
  ghcFlow: false,
};

const environmentFlags = {
  development: Object.assign({}, defaultFlags, {
    allOrdersTypes: true,
    hhgFlow: true,
    ghcFlow: true,
  }),

  test: Object.assign({}, defaultFlags),

  experimental: Object.assign({}, defaultFlags, {
    allOrdersTypes: true,
    hhgFlow: true,
    ghcFlow: true,
  }),

  staging: Object.assign({}, defaultFlags, {
    allOrdersTypes: true,
    hhgFlow: true,
    ghcFlow: true,
  }),

  production: Object.assign({}, defaultFlags, {
    sitPanel: false,
  }),
};

export function flagsFromURL(search) {
  const params = queryString.parse(search);
  let flags = {};

  forEach(params, function (value, key) {
    let [prefix, name] = key.split(':');
    if (prefix === 'flag' && name.length > 0) {
      if (validateFlag(name)) {
        // name is validated by the previous line
        //  security/detect-object-injection
        flags[name] = value === 'true';
      }
    }
  });
  return flags;
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
    case 'office.staging.move.mil':
    case 'my.staging.move.mil':
    case 'admin.staging.move.mil':
      return 'staging';
    case 'my.experimental.move.mil':
    case 'office.experimental.move.mil':
    case 'admin.experimental.move.mil':
      return 'experimental';
    default:
      return 'development';
  }
}

function validateFlag(name) {
  // Warn if the value being fetched was never set.
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
  //  security/detect-object-injection
  return Object.assign({}, environmentFlags[env], flagsFromURL(search));
}

export const createModifiedSchemaForOrdersTypesFlag = (schema) => {
  const ordersTypeSchema = Object.assign({}, schema.properties.orders_type, {
    enum: [schema.properties.orders_type.enum[0]],
  });
  const properties = Object.assign({}, schema.properties, { orders_type: ordersTypeSchema });
  const modifiedSchema = Object.assign({}, schema, {
    properties: properties,
  });

  return modifiedSchema;
};
