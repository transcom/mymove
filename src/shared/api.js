import { get, uniqueId } from 'lodash';
import Swagger from 'swagger-client';
import { normalize } from 'normalizr';

import store from 'shared/store';
import { addEntities } from 'shared/Entities/actions'
let client = null;
let publicClient = null;

export function request(label, opName, params) {
  return async function(dispatch, getState, { schema }) {
    const client = await getClient();
    const operation = get(client, 'apis.' + opName);

    if (!operation) {
      throw new Error(`Operation '${opName}' does not exist!`);
    }

    const id = uniqueId('req_');
    const requestLog = {
      id,
      opName,
      params,
      start: new Date(),
    };
    store.dispatch({ type: `@@swagger/${opName}/START`, label, request: requestLog });

    const request = operation(params)
      .then(response => {
        const updatedRequestLog = Object.assign({}, requestLog, {
          ok: response.ok,
          end: new Date(),
        });
        dispatch({ type: `@@swagger/${opName}/SUCCESS`, response, request: updatedRequestLog });
        const data = normalize(response.body, schema.shipment);
        dispatch(addEntities(data.entities));
        return response;
      })
      .catch(response => {
        console.error(`Operation ${opName} failed: ${response} (${response.status})`);
        const updatedRequestLog = Object.assign({}, requestLog, {
          ok: false,
          end: new Date(),
          response,
        });
        dispatch({ type: `@@swagger/${opName}/FAILURE`, response, request: updatedRequestLog });
        throw response;
      });
    return request;
  };
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
