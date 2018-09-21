import { getClient, checkResponse } from 'shared/Swagger/api';

export async function IssuesIndex() {
  const client = await getClient();
  const response = await client.apis.issues.indexIssues();
  checkResponse(response, 'failed to load issues index due to server error');
  return response.body;
}
