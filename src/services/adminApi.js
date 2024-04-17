import Swagger from 'swagger-client';

import { makeSwaggerRequest, requestInterceptor, responseInterceptor, makeSwaggerRequestRaw } from './swaggerRequest';

let adminClient = null;

// setting up the same config from Swagger/api.js
export async function getAdminClient() {
  if (!adminClient) {
    adminClient = await Swagger({
      url: '/admin/v1/swagger.yaml',
      requestInterceptor,
      responseInterceptor,
    });
  }
  return adminClient;
}

export async function makeAdminRequest(operationPath, params = {}, options = {}) {
  const client = await getAdminClient();
  return makeSwaggerRequest(client, operationPath, params, options);
}

export async function makeAdminRequestRaw(operationPath, params = {}) {
  const client = await getAdminClient();
  return makeSwaggerRequestRaw(client, operationPath, params);
}

export async function updateRequestedOfficeUser(officeUserId, body) {
  const operationPath = 'Requested office users.updateRequestedOfficeUser';

  return makeAdminRequest(
    operationPath,
    {
      officeUserId,
      body,
    },
    { normalize: false },
  );
}
