import Swagger from 'swagger-client';
let client = null;

async function ensureClientIsLoaded() {
  if (!client) {
    client = await Swagger('/api/v1/swagger.yaml');
  }
}

function checkResponse(response, errorMessage) {
  if (!response.ok) {
    throw new Error(`${errorMessage}: ${response.url}: ${response.statusText}`);
  }
}

export async function GetSpec() {
  await ensureClientIsLoaded();
  return client.spec;
}
// because these are async functions,
// they return a promise that is resolved with the return value
// if there is an error, the promise is rejected with that error
export async function IssuesIndex() {
  await ensureClientIsLoaded();
  const response = await client.apis.issues.indexIssues();
  checkResponse(response, 'failed to load issues index due to server error');
  return response.body;
}

export async function CreateIssue(issueBody) {
  await ensureClientIsLoaded();
  const response = await client.apis.issues.createIssue({
    createIssuePayload: { description: issueBody },
  });
  checkResponse(response, 'failed to create issue due to server error');
}
export async function CreateForm1299(formData) {
  await ensureClientIsLoaded();
  const response = await client.apis.form1299s.createForm1299({
    createForm1299Payload: formData,
  });
  checkResponse(response, 'failed to create form 1299 due to server error');
  //todo: return uuid?
}

export async function ShipmentsIndex() {
  await ensureClientIsLoaded();
  const response = await client.apis.shipments.indexShipments();
  checkResponse(response, 'failed to load shipments index due to server error');
  return response.body;
}
