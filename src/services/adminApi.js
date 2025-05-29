/* istanbul ignore file */
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

export async function getTransportationOfficeByID(officeId) {
  const operationPath = 'Transportation offices.getOfficeById';
  return makeAdminRequest(operationPath, { officeId }, { normalize: false });
}

export async function deleteOfficeUser(officeUserId) {
  const operationPath = 'Office users.deleteOfficeUser';

  return makeAdminRequest(
    operationPath,
    {
      officeUserId,
    },
    { normalize: false },
  );
}

export async function updateOfficeUser(officeUserId, officeUser) {
  const operationPath = 'Office users.updateOfficeUser';

  return makeAdminRequest(
    operationPath,
    {
      officeUserId,
      officeUser,
    },
    { normalize: false },
  );
}

export async function getRolesPrivileges() {
  const operationPath = 'Office users.getRolesPrivileges';
  return makeAdminRequest(operationPath, {}, { normalize: false });
}

export async function deleteUser(userId) {
  const operationPath = 'Users.deleteUser';

  return makeAdminRequest(
    operationPath,
    {
      userId,
    },
    { normalize: false },
  );
}

export async function updateUser(userId, user) {
  const operationPath = 'Users.updateUser';

  return makeAdminRequest(
    operationPath,
    {
      userId,
      User: user,
    },
    { normalize: false },
  );
}
