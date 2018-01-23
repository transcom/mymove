import Swagger from 'swagger-client';
let client = null;

async function ensureClientIsLoaded() {
  if (!client) {
    client = await Swagger('api/v1/swagger.yaml');
  }
}

export async function IssuesIndex() {
  await ensureClientIsLoaded();
  const result = await client.apis.default.indexIssues();
  if (result.ok) return result.body;
  else {
    throw new Error('failed to load issues index');
  }
}

export async function CreateIssue(issueBody) {
  await ensureClientIsLoaded();
  const result = await client.apis.default.createIssue({
    issue: { body: issueBody },
  });
  if (!result.ok) throw new Error('failed to createIssue');
}
