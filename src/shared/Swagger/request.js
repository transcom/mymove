import { get, uniqueId } from 'lodash';
import { getClient } from './api';
import { normalize } from 'normalizr';

import store from 'shared/store';

// Given a schema path (e.g. shipments.getShipment), return the
// route's definition from the Swagger spec
function findMatchingRoute(paths, operationPath) {
  const [tagName, operationId] = operationPath.split('.');

  let routeDefinition;
  Object.values(getSpec().paths).some(function(path) {
    return Object.values(path).some(function(route) {
      if (route.operationId === operationId && route.tags[0] === tagName) {
        routeDefinition = route;
        return true;
      }
      return false;
    });
  });
  return routeDefinition;
}

// Given a route definition and a status code, return the lowercased
// name for the defined return type. For example, a 200 response to
// shipments.getShipment should be '$$ref/definitions/Shipment', for
// which this function would return 'shipment'.
//
// This key can be used to determine what key to find the object's
// definition in within our normalizr schema.
function successfulReturnType(routeDefinition, status) {
  const response = routeDefinition.responses[status];
  if (!response) {
    console.error(
      `No response found for operation ${
        routeDefinition.operationId
      } with status ${status}`,
    );
    return;
  }
  return response.schema['$$ref'].split('/')[2].toLowerCase();
}

// Call an operation defined in the Swagger API, dispatching
// actions as its state changes.
export function swaggerRequest(operationPath, params, options = {}) {
  return async function(dispatch, getState, { schema }) {
    const client = await getClient();
    const operation = get(client, 'apis.' + operationPath);

    if (!operation) {
      throw new Error(`Operation '${operationPath}' does not exist!`);
    }

    const id = uniqueId('req_');
    const label = options.label || id;
    const requestLog = {
      id,
      operationPath,
      params,
      start: new Date(),
      isLoading: true,
    };
    store.dispatch({
      type: `@@swagger/${operationPath}/START`,
      label,
      request: requestLog,
    });

    let request;

    try {
      request = operation(params);
    } catch (error) {
      console.error(`Operation ${operationPath} failed: ${error}`);
      const updatedRequestLog = Object.assign({}, requestLog, {
        ok: false,
        end: new Date(),
        isLoading: false,
        error,
      });
      dispatch({
        type: `@@swagger/${operationPath}/ERROR`,
        error,
        request: updatedRequestLog,
      });
    }

    request
      .then(response => {
        const updatedRequestLog = Object.assign({}, requestLog, {
          ok: response.ok,
          end: new Date(),
          isLoading: false,
        });

        const action = {
          type: `@@swagger/${operationPath}/SUCCESS`,
          request: updatedRequestLog,
          response,
        };

        const routeDefinition = findMatchingRoute(
          client.spec.paths,
          operationPath,
        );
        const schemaKey = successfulReturnType(
          routeDefinition,
          response.status,
        );

        const payloadSchema = schema[schemaKey];
        if (payloadSchema) {
          action.entities = normalizePayload(
            response.body,
            payloadSchema,
          ).entities;
        } else {
          console.warn(`Could not find a schema for ${schemaKey}`);
        }
        dispatch(action);
        return response;
      })
      .catch(response => {
        console.error(
          `Operation ${operationPath} failed: ${response} (${response.status})`,
        );
        const updatedRequestLog = Object.assign({}, requestLog, {
          ok: false,
          end: new Date(),
          response,
          isLoading: false,
        });
        dispatch({
          type: `@@swagger/${operationPath}/FAILURE`,
          response,
          request: updatedRequestLog,
        });
        throw response;
      });
    return request;
  };
}

function normalizePayload(body, schema) {
  return normalize(body, schema);
}
