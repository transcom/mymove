import Swagger from 'swagger-client';

import { makeSwaggerRequest, requestInterceptor, responseInterceptor } from './swaggerRequest';

let primeSimulatorClient = null;

// setting up the same config from Swagger/api.js
export async function getPrimeSimulatorClient() {
  if (!primeSimulatorClient) {
    primeSimulatorClient = await Swagger({
      url: '/prime/v1/swagger.yaml',
      requestInterceptor,
      responseInterceptor,
    });
  }
  return primeSimulatorClient;
}

export async function makePrimeSimulatorRequest(operationPath, params = {}, options = {}) {
  const client = await getPrimeSimulatorClient();
  return makeSwaggerRequest(client, operationPath, params, options);
}

export async function getPrimeSimulatorAvailableMoves() {
  const operationPath = 'moveTaskOrder.listMoves';
  return makePrimeSimulatorRequest(operationPath, {}, { schemaKey: 'listMoves', normalize: false });
}

export async function getPrimeSimulatorMove(key, locator) {
  return makePrimeSimulatorRequest('moveTaskOrder.getMoveTaskOrder', { moveID: locator }, { normalize: true });
}
