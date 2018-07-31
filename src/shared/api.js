import {get, uniqueId} from 'lodash';
import Swagger from 'swagger-client';

import store from 'shared/store';

let client = null;

window.request = function(label, opName, params) {
  const operation = get(client, 'apis.' + opName);

  if (!operation) {
    throw(new Error(`Operation '${opName}' does not exist!`));
  }

  const id = uniqueId('req_');

  store.dispatch({type: `@@swagger/${opName}/START`, label, id});

  return operation(params)
    .then(payload => store.dispatch({type: `@@swagger/${opName}/SUCCESS`, id, payload}))
    .catch(error => store.dispatch({type: `@@swagger/${opName}/FAILURE`, id, error}));
}

export async function getClient() {
  if (!client) {
    client = await Swagger({
      url: '/internal/swagger.yaml',
      requestInterceptor: req => {
        console.debug(req);
      }
    });
  }
  return client;
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
