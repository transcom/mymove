import Swagger from 'swagger-client';
import * as Cookies from 'js-cookie';
import { setLoggedOut } from 'shared/Data/users';

let client = null;
let publicClient = null;

export const requestInterceptor = req => {
  if (!req.loadSpec) {
    const token = Cookies.get('masked_gorilla_csrf');
    if (token) {
      req.headers['X-CSRF-Token'] = token;
    } else {
      console.warn('Unable to retrieve CSRF Token from cookie');
    }
  }
  return req;
};

export const responseInterceptor = res => {
  if (res.status === 401) {
    setLoggedOut();
  }
  return res;
};

export async function getClient() {
  if (!client) {
    client = await Swagger({
      url: '/internal/swagger.yaml',
      requestInterceptor: requestInterceptor,
      responseInterceptor: responseInterceptor,
    });
  }
  return client;
}

export async function getPublicClient() {
  if (!publicClient) {
    publicClient = await Swagger({
      url: '/api/v1/swagger.yaml',
      requestInterceptor: requestInterceptor,
      responseInterceptor: responseInterceptor,
    });
  }
  return publicClient;
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
    throw new Error(`${errorMessage}: ${response.url}: ${response.statusText}`);
  }
}
