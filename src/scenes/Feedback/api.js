import { ensureClientIsLoaded, checkResponse } from 'shared/api';

export async function CreateIssue(issueBody) {
  const client = await ensureClientIsLoaded();
  const response = await client.apis.issues.createIssue({
    createIssuePayload: { description: issueBody },
  });
  checkResponse(response, 'failed to create issue due to server error');
}
