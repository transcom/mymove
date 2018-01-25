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
  const response = await client.apis.default.indexIssues();
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
  const response = await client.apis.default.createIssue({
    issue: { body: issueBody },
  });
  if (!response.ok)
    new Error(
      `failed to create issue due to server error:
      ${response.url}: ${response.statusText}`,
    );
}
