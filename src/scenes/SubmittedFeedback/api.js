import { ensureClientIsLoaded, checkResponse } from 'shared/api';

export async function IssuesIndex() {
  const client = await ensureClientIsLoaded();
  const response = await client.apis.issues.indexIssues();
  checkResponse(response, 'failed to load issues index due to server error');
  return response.body;
}
