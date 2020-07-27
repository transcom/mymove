import Swagger from 'swagger-client';
import * as Cookies from 'js-cookie';
import { get } from 'lodash';
import { normalize } from 'normalizr';

import * as schema from 'shared/Entities/schema';

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

let ghcClient = null;

// setting up the same config from Swagger/api.js
export async function getGHCClient() {
  if (!ghcClient) {
    ghcClient = await Swagger({
      url: '/ghc/v1/swagger.yaml',
      requestInterceptor,
    });
  }
  return ghcClient;
}

// this is what makes the API call instead of SwaggerRequest (but very similar)
export async function getPaymentRequestList() {
  const operationPath = 'paymentRequests.listPaymentRequests';
  const client = await getGHCClient();
  const operation = get(client, `apis.${operationPath}`);

  if (!operation) {
    // eslint-disable-next-line no-console
    console.log('Operation does not exist', operationPath);
  }

  let request;
  try {
    request = operation();
  } catch (e) {
    // eslint-disable-next-line no-console
    console.error('operation failed', e);
  }

  return request
    .then((response) => {
      const payloadSchema = schema.paymentRequests;
      return normalize(response.body, payloadSchema).entities;
    })
    .catch((response) => {
      // eslint-disable-next-line no-console
      console.log('response failed', response);
      return Promise.reject();
    });
}
