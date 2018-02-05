import Swagger from 'swagger-client';
let client = null;

async function ensureClientIsLoaded() {
  if (!client) {
    client = await Swagger('api/v1/swagger.yaml');
  }
}

// because these are async functions,
// they return a promise that is resolved with the return value
// if there is an error, the promise is rejected with that error
export async function IssuesIndex() {
  await ensureClientIsLoaded();
  const response = await client.apis.issues.indexIssues();
  if (!response.ok) {
    throw new Error(
      `failed to load issues index due to server error:
      ${response.url}: ${response.statusText}`,
    );
  }
  return response.body;
}

export async function CreateIssue(issueBody) {
  await ensureClientIsLoaded();
  const response = await client.apis.issues.createIssue({
    createIssuePayload: { description: issueBody },
  });
  if (!response.ok)
    new Error(
      `failed to create issue due to server error:
      ${response.url}: ${response.statusText}`,
    );
}

export async function ShipmentIndex() {
  await ensureClientIsLoaded();
  const response = await client.apis.shipments.indexIssues();
  if (!response.ok) {
    throw new Error(
      `failed to load shipments index due to server error:
      ${response.url}: ${response.statusText}`,
    );
  }
  return response.body;
}
