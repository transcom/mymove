import Swagger from 'swagger-client';
let client = null;

async function ensureClientIsLoaded() {
  if (!client) {
    client = await Swagger('api/v1/swagger.yaml');
  }
}

export async function IssuesIndex() {
  await ensureClientIsLoaded();
  return client.apis.default.indexIssues();
}

export async function CreateIssue(issueBody) {
  await ensureClientIsLoaded();
  return client.apis.default.createIssue({ issue: { body: issueBody } });
}
