import { filter, get, last, sortBy } from 'lodash';

// Return a convenient object that contains commonly needed info about
// the requests for a label
export function getRequestStatus(state, label) {
  return {
    error: getLastError(state, label),
    isLoading: getLastRequestIsLoading(state, label),
  };
}

// Get the last request for a provided label
export function getLastRequest(state, label) {
  const requests = filter(state.requests.byID, function(value, key) {
    return value.label === label;
  });
  const sorted = sortBy(requests, ['start']);
  return last(sorted);
}

// Return if the last request for
export function getLastRequestIsLoading(state, label) {
  const last = getLastRequest(state, label);
  if (last) {
    return last.isLoading;
  } else {
    return false;
  }
}

// Return the last error for a given label
export function getLastError(state, label) {
  // eslint-disable-next-line security/detect-object-injection
  return state.requests.lastErrors[label];
}

// Return the Swagger definition for the provided name
export function getInternalSwaggerDefinition(state, name) {
  return get(state, `swaggerInternal.spec.definitions.${name}`, {});
}
