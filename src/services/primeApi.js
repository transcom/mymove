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

export async function getPrimeSimulatorAvailableMoves(
  key,
  { sort, order, filters = [], currentPage = 1, currentPageSize = 20 },
) {
  const operationPath = 'moveTaskOrder.listMoves';
  const paramFilters = {};
  filters.forEach((filter) => {
    paramFilters[`${filter.id}`] = filter.value;
  });
  return makePrimeSimulatorRequest(
    operationPath,
    {
      sort,
      order,
      page: currentPage,
      perPage: currentPageSize,
      ...paramFilters,
    },
    { schemaKey: 'listMoves', normalize: false },
  );
}
