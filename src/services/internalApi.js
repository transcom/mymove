import Swagger from 'swagger-client';
import * as Cookies from 'js-cookie';

import { makeSwaggerRequest } from './swaggerRequest';

// setting up the same config from Swagger/api.js
const requestInterceptor = (req) => {
  if (!req.loadSpec) {
    const token = Cookies.get('masked_gorilla_csrf');
    if (token) {
      req.headers['X-CSRF-Token'] = token;
    } else {
      // eslint-disable-next-line no-console
      console.warn('Unable to retrieve CSRF Token from cookie');
    }
  }
  return req;
};

let internalClient = null;

// setting up the same config from Swagger/api.js
export async function getInternalClient() {
  if (!internalClient) {
    internalClient = await Swagger({
      url: '/internal/swagger.yaml',
      requestInterceptor,
    });
  }
  return internalClient;
}

export async function makeInternalRequest(operationPath, params = {}, options = {}) {
  const client = await getInternalClient();
  return makeSwaggerRequest(client, operationPath, params, options);
}

export async function getMoveOrder(key, ordersId) {
  return makeInternalRequest('orders.showOrders', { ordersId });
}

export async function updateOrders({ ordersId, body }) {
  return makeInternalRequest('orders.updateOrders', { ordersId, updateOrders: body });
}
