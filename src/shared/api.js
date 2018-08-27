import { some, get, uniqueId } from 'lodash';
import Swagger from 'swagger-client';
import { normalize } from 'normalizr';

import store from 'shared/store';
let client = null;
let publicClient = null;

// Given a schema path (e.g. shipments.getShipment), return the
// route's definition from the Swagger spec
function findMatchingRoute(paths, operationPath) {
  const [tagName, operationId] = operationPath.split('.');

  let routeDefinition;
  Object.values(client.spec.paths).some(function(path) {
    return Object.values(path).some(function(route) {
      if (route.operationId === operationId && route.tags[0] === tagName) {
        routeDefinition = route;
        return true;
      }
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

export async function getClient() {
  if (!client) {
    client = await Swagger({ url: '/internal/swagger.yaml' });
  }
  return client;
}

export async function getPublicClient() {
  if (!publicClient) {
    publicClient = await Swagger('/api/v1/swagger.yaml');
  }
  return publicClient;
}

export function checkResponse(response, errorMessage) {
  if (!response.ok) {
    throw new Error(`${errorMessage}: ${response.url}: ${response.statusText}`);
  }
}

export async function CreateUpload(fileUpload, documentId) {
  const client = await getClient();
  const response = await client.apis.uploads.createUpload({
    file: fileUpload,
    documentId,
  });
  checkResponse(response, 'failed to upload file due to server error');
  return response.body;
}

export async function DeleteUpload(uploadId) {
  const client = await getClient();
  const response = await client.apis.uploads.deleteUpload({
    uploadId,
  });
  checkResponse(response, 'failed to delete file due to server error');
  return response.body;
}

export async function DeleteUploads(uploadIds) {
  const client = await getClient();
  const response = await client.apis.uploads.deleteUploads({
    uploadIds,
  });
  checkResponse(response, 'failed to delete files due to server error');
  return response.body;
}

export async function CreateDocument(name, serviceMemberId) {
  const client = await getClient();
  const response = await client.apis.documents.createDocument({
    documentPayload: {
      name: name,
      service_member_id: serviceMemberId,
    },
  });
  checkResponse(response, 'failed to create document due to server error');
  return response.body;
}
