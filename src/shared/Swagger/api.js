import Swagger from 'swagger-client';
import Cookies from 'js-cookie';

import { getInternalClient } from 'services/internalApi';
import { milmoveLogger } from 'utils/milmoveLog';

let publicClient = null;
let ghcClient = null;
let adminClient = null;

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

export async function getClient() {
  return await getInternalClient();
}

export async function getPublicClient() {
  if (!publicClient) {
    publicClient = await Swagger({
      url: '/api/v1/swagger.yaml',
      requestInterceptor,
    });
  }
  return publicClient;
}

export async function getGHCClient() {
  if (!ghcClient) {
    ghcClient = await Swagger({
      url: '/ghc/v1/swagger.yaml',
      requestInterceptor,
    });
  }
  return ghcClient;
}

export async function getAdminClient() {
  if (!adminClient) {
    adminClient = await Swagger({
      url: '/admin/v1/swagger.yaml',
      requestInterceptor,
    });
  }
  return adminClient;
}

export async function getSpec() {
  const client = await getClient();
  return client.spec;
}

export async function getPublicSpec() {
  const client = await getPublicClient();
  return client.spec;
}

// Used by pre-swaggerRequest code to verify responses
export function checkResponse(response, errorMessage) {
  if (!response.ok) {
    const err = new Error(`${errorMessage}: ${response.url}: ${response.statusText}`);
    err.status = response.status;
    throw err;
  }
}
