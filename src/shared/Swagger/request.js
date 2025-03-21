import { get, some, uniqueId } from 'lodash';
import { normalize } from 'normalizr';
import { milmoveLogger } from 'utils/milmoveLog';

// Given a schema path (e.g. shipments.getShipment), return the
// route's definition from the Swagger spec
function findMatchingRoute(paths, operationPath) {
  const [tagName, operationId] = operationPath.split('.');

  let routeDefinition;
  Object.values(paths).some(function (path) {
    return some(path, function (route, method) {
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
    return undefined;
  }
  return toCamelCase(schemaKey);
}

// Call an operation defined in the Swagger API, dispatching
// actions as its state changes.
export function swaggerRequest(getClient, operationPath, params, options = {}) {
  return async function (dispatch, getState, { schema }) {
    const client = await getClient();
    const state = await getState();
    const operation = get(client, `apis.${operationPath}`);

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
      label,
    };
    dispatch({
      type: `@@swagger/${operationPath}/START`,
      request: requestLog,
      label,
    });

    let request;

    try {
      request = operation(params);
    } catch (error) {
      milmoveLogger.error(`Operation ${operationPath} failed: ${error}`);
      const updatedRequestLog = { ...requestLog, ok: false, end: new Date(), isLoading: false, error };
      dispatch({
        type: `@@swagger/${operationPath}/ERROR`,
        request: updatedRequestLog,
        error,
        label,
      });
      return Promise.reject(error);
    }

    return request
      .then((response) => {
        const updatedRequestLog = { ...requestLog, ok: response.ok, end: new Date(), isLoading: false };
        const routeDefinition = findMatchingRoute(client.spec.paths, operationPath);
        if (!routeDefinition) {
          throw new Error(`Could not find routeDefinition for ${operationPath}`);
        }

        const action = {
          type: `@@swagger/${operationPath}/SUCCESS`,
          request: updatedRequestLog,
          method: routeDefinition.method,
          response,
          label,
        };

        let { schemaKey } = options;
        if (!schemaKey) {
          schemaKey = successfulReturnType(routeDefinition, response.status);
        }
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

        const payloadSchema = schema[schemaKey];
        if (!payloadSchema) {
          throw new Error(`Could not find a schema for ${schemaKey}`);
        }
        if (options.deleteId) {
          const oldEntity = state.entities[schemaKey][options.deleteId];
          action.entities = normalize([oldEntity], payloadSchema).entities;
        } else {
          action.entities = normalize(response.body, payloadSchema).entities;
        }
        dispatch(action);
        return action;
      })
      .catch((responseError) => {
        const traceId = responseError?.response?.headers['x-milmove-trace-id'] || 'unknown-milmove-trace-id';
        milmoveLogger.error(
          `Operation ${operationPath} failed: ${responseError} (${responseError.status})`,
          `milmove_trace_id: ${traceId}`,
        );
        const updatedRequestLog = {
          ...requestLog,
          ok: false,
          end: new Date(),
          response: responseError.response,
          isLoading: false,
        };
        const action = {
          type: `@@swagger/${operationPath}/FAILURE`,
          request: updatedRequestLog,
          response: responseError.response,
          label,
        };
        dispatch(action);
        return Promise.reject(action);
      });
  };
}

export function resetRequests() {
  return {
    type: '@@swagger/RESET',
  };
}
