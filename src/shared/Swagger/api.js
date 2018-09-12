import Swagger from 'swagger-client';

let client = null;
let publicClient = null;

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
