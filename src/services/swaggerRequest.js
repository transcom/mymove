/* eslint-disable import/prefer-default-export */
import { get } from 'lodash';
import { normalize } from 'normalizr';
import Cookies from 'js-cookie';

import * as schema from 'shared/Entities/schema';
import { interceptInjection } from 'store/interceptor/injectionMiddleware';
import { interceptResponse } from 'store/interceptor/actions';
import { milmoveLogger } from 'utils/milmoveLog';
import { setMarkerIOMilmoveTraceID } from 'components/ThirdParty/MarkerIO';

// setting up the same config from Swagger/api.js
export const requestInterceptor = (req) => {
  if (!req.loadSpec) {
    const token = Cookies.get('masked_gorilla_csrf');
    if (token) {
      req.headers['X-CSRF-Token'] = token;
    } else {
      milmoveLogger.warn('Unable to retrieve CSRF Token from cookie');
    }
  }
  return req;
};

export const responseInterceptor = (res) => {
  switch (res.status) {
    case 500: {
      interceptInjection(interceptResponse(true, res.headers['x-milmove-trace-id']));
      break;
    }

    default: {
      interceptInjection(interceptResponse(false));
    }
  }

  if (!res.ok) {
    setMarkerIOMilmoveTraceID(res.headers['x-milmove-trace-id']);
  }
  return res;
};

/**
 * This is a new Swagger request function that does not rely on Redux
 */

// Given a schema path (e.g. shipments.getShipment), return the
// route's definition from the Swagger spec
function findMatchingRoute(paths, operationPath) {
  const [tagName, operationId] = operationPath.split('.');

  let routeDefinition;
  Object.values(paths).some((path) => {
    return Object.values(path).some((route, method) => {
      if (route.operationId === operationId && route.tags[0] === tagName) {
        routeDefinition = route;
        routeDefinition.method = method;
        return true;
      }
      return false;
    });
  });
  return routeDefinition;
}

// assumes str passed in is title case (ex. SomeModelName => someModelName)
const toCamelCase = (str) => str[0].toLowerCase() + str.slice(1);

// Given a route definition and a status code, return the lowercased
// name for the defined return type. For example, a 200 response to
// shipments.getShipment should be '$$ref/definitions/Shipment', for
// which this function would return 'shipment'.
//
// This key can be used to determine what key to find the object's
// definition in within our normalizr schema.
function successfulReturnType(routeDefinition, status) {
  const response = routeDefinition.responses[status];
  const schemaKey = response.schema.$$ref.split('/').pop();
  if (!response) {
    milmoveLogger.error(`No response found for operation ${routeDefinition.operationId} with status ${status}`);
    return null;
  }

  return toCamelCase(schemaKey);
}

export function normalizeResponse(data, schemaKey) {
  const responseSchema = schema[`${schemaKey}`];
  if (!responseSchema) {
    throw new Error(`Could not find a schema for ${schemaKey}`);
  }

  return normalize(data, responseSchema).entities;
}

export async function makeSwaggerRequest(client, operationPath, params = {}, options = { normalize: true }) {
  const operation = get(client, `apis.${operationPath}`);
  if (!operation) {
    throw new Error(`Operation '${operationPath}' does not exist!`);
  }

  let request;
  try {
    request = operation(params);
  } catch (e) {
    milmoveLogger.error(`Operation ${operationPath} failed: ${e}`);
    return Promise.reject(e);
  }

  return request
    .then((response) => {
      const normalizeData = options.normalize !== undefined ? options.normalize : true;
      // Normalize the data (defaults to true)
      if (normalizeData) {
        /* TODO - deprecrate the below & require an explicit schemaKey parameter */
        const routeDefinition = findMatchingRoute(client.spec.paths, operationPath);
        if (!routeDefinition) {
          throw new Error(`Could not find routeDefinition for ${operationPath}`);
        }

        let schemaKey = options.schemaKey || successfulReturnType(routeDefinition, response.status);
        if (!schemaKey) {
          throw new Error(`Could not find schemaKey for ${operationPath} status ${response.status}`);
        }

        if (schemaKey.indexOf('Payload') !== -1) {
          const newSchemaKey = schemaKey.replace('Payload', '');
          milmoveLogger.warn(
            `Using 'Payload' as a response type prefix is deprecated. Please rename ${schemaKey} to ${newSchemaKey}`,
          );
          schemaKey = newSchemaKey;
        }
        return normalizeResponse(response.body, schemaKey);
      }

      // Otherwise, return raw response body
      return response.body;
    })
    .catch((responseError) => {
      const traceId = responseError?.response?.headers['x-milmove-trace-id'] || 'unknown-milmove-trace-id';
      milmoveLogger.error(
        `Operation ${operationPath} failed: ${responseError} (${responseError.status})`,
        `milmove_trace_id: ${traceId}`,
      );
      return Promise.reject(responseError);
    });
}

export async function makeSwaggerRequestRaw(client, operationPath, params = {}) {
  const operation = get(client, `apis.${operationPath}`);
  if (!operation) {
    throw new Error(`Operation '${operationPath}' does not exist!`);
  }

  let request;
  try {
    request = operation(params);
  } catch (e) {
    milmoveLogger.error(`Operation ${operationPath} failed: ${e}`);
    return Promise.reject(e);
  }

  return request
    .then((response) => {
      // no normalizaton ..etc...just return entire response
      return response;
    })
    .catch((responseError) => {
      const traceId = responseError?.response?.headers['x-milmove-trace-id'] || 'unknown-milmove-trace-id';
      milmoveLogger.error(
        `Operation ${operationPath} failed: ${responseError} (${responseError.status})`,
        `milmove_trace_id: ${traceId}`,
      );
      return Promise.reject(responseError);
    });
}
